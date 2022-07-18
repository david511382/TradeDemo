package api

import (
	"sync"
	tradeBll "zerologix-homework/src/bll/trade"
	"zerologix-homework/src/server/reqs"
	"zerologix-homework/src/server/resp"
)

type Order struct {
	tradeLogic tradeBll.ITradeLogic
}

func NewOrder(
	tradeLogic tradeBll.ITradeLogic,
) *Order {
	return &Order{
		tradeLogic: tradeLogic,
	}
}

func (l Order) PostTest(req *reqs.OrderPostTest) (
	result resp.OrderPostTest,
	resultErr error,
) {
	runTimes := req.RunTimes
	r := req.OrderPostBuy
	wg := sync.WaitGroup{}
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.PostBuy(&r)
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.PostSell(&reqs.OrderPostSell{
				OrderPostBuy: r,
			})
		}()
	}

	wg.Wait()
	result.Message = "OK"

	return
}

func (l Order) PostBuy(req *reqs.OrderPostBuy) (
	result resp.OrderPostBuy,
	resultErr error,
) {
	if err := l.tradeLogic.Match(tradeBll.Order{
		OrderType: tradeBll.ORDER_TYPE_BUY,
		Quantity:  req.Quantity,
		Price:     req.Price,
	}); err != nil {
		resultErr = err
		result.Message = "Error"
		return
	}

	result.Message = "OK"

	return
}

func (l Order) PostSell(req *reqs.OrderPostSell) (
	result resp.OrderPostSell,
	resultErr error,
) {
	if err := l.tradeLogic.Match(tradeBll.Order{
		OrderType: tradeBll.ORDER_TYPE_SELL,
		Quantity:  req.Quantity,
		Price:     req.Price,
	}); err != nil {
		resultErr = err
		result.Message = "Error"
		return
	}

	result.Message = "OK"

	return
}
