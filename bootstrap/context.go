package bootstrap

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Context struct {
	ctx        context.Context
	cancel     context.CancelFunc
	closeSigCh chan os.Signal
	notifyCh   []chan os.Signal
	sigOnce    *sync.Once
	wg         *sync.WaitGroup
}

func NewContext() *Context {
	ctx, cancel := context.WithCancel(context.Background())
	c := &Context{
		ctx,
		cancel,
		make(chan os.Signal, 1),
		make([]chan os.Signal, 0),
		new(sync.Once),
		&sync.WaitGroup{},
	}
	c.setupGracefulSignal()
	return c
}

// 發送關閉信號
func (c *Context) Shutdown() {
	if c == nil {
		return
	}
	c.sigOnce.Do(c.cancel)
}

// 關閉
func (c *Context) shutdownBy(s os.Signal) {
	for _, ch := range c.notifyCh {
		ch <- s
	}
	c.sigOnce.Do(c.cancel)
}

// 設定優雅關閉的信號
func (c *Context) setupGracefulSignal() {
	signal.Notify(
		c.closeSigCh,

		// 中斷，當使用者從鍵盤按ctrl+c鍵
		syscall.SIGINT,

		// 軟體終止（software? termination）
		// k8s 關閉pod會傳這個信號
		// 傳送後會等 terminationGracePeriodSeconds，預設30s
		// 等待 terminationGracePeriodSeconds 時間後會傳 SIGKILL
		syscall.SIGTERM,
	)
	go func() {
		select {
		case <-c.ctx.Done():
			c.shutdownBy(syscall.SIGTERM)
		case s := <-c.closeSigCh:
			c.shutdownBy(s)
		}
	}()
}

// GracefulDown 監聽關閉
func (c *Context) GracefulDown(cancel func()) {
	go c.GracefulDownBlock(func(s os.Signal) {
		if cancel != nil {
			cancel()
		}
	})
}

// GracefulDown 監聽關閉信號
func (c *Context) GracefulDownSignal(cancel func(os.Signal)) {
	go c.GracefulDownBlock(cancel)
}

// GracefulDown 監聽關閉信號並阻塞
func (c *Context) GracefulDownBlock(cancel func(os.Signal)) {
	ch := make(chan os.Signal, 1)
	c.notifyCh = append(c.notifyCh, ch)

	c.wg.Add(1)
	s := <-ch
	if cancel != nil {
		cancel(s)
	}
	c.wg.Done()

	c.wg.Wait()
}
