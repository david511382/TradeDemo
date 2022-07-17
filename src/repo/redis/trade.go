package redis

import (
	"zerologix-homework/bootstrap"
	"zerologix-homework/src/repo/redis/schema/trade"
)

func Trade() *trade.Database {
	lock.RLock()
	isNoValue := tradeDb == nil
	lock.RUnlock()
	if isNoValue {
		lock.Lock()
		defer lock.Unlock()
		if tradeDb == nil {
			repo := TradeRedisCfgRepo{}
			keyRoot := ""
			if cfg, err := bootstrap.Get(); err == nil {
				_, key := repo.Get(cfg)
				keyRoot = key
			}
			tradeDb = trade.NewDatabase(
				getConnect(func(cfg *bootstrap.Config) bootstrap.Db {
					c, _ := repo.Get(cfg)
					return c
				}),
				keyRoot,
			)
		}
	}
	copy := *tradeDb
	return &copy
}

type TradeRedisCfgRepo struct{}

func (TradeRedisCfgRepo) Get(cfg *bootstrap.Config) (repoCfg bootstrap.Db, keyRoot string) {
	repoCfg = cfg.RedisConfig.Trade
	keyRoot = cfg.RedisConfig.RedisTradeKeyRoot + ":"
	return
}

func (TradeRedisCfgRepo) Set(cfg *bootstrap.Config, testName string) {
	cfg.RedisConfig.RedisTradeKeyRoot = testName
}
