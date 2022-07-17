package redis

import (
	"sync"
	"time"
	"zerologix-homework/bootstrap"
	"zerologix-homework/src/repo/redis/conn"
	"zerologix-homework/src/repo/redis/schema/trade"

	"github.com/go-redis/redis"
)

var (
	tradeDb *trade.Database
	lock    sync.RWMutex
)

func getConnect(configSelector func(cfg *bootstrap.Config) bootstrap.Db) func() (master, slave *redis.Client, resultErr error) {
	return GetConnectFn(
		bootstrap.Get,
		configSelector,
	)
}

func GetConnectFn(
	configGetterFn func() (*bootstrap.Config, error),
	configSelector func(cfg *bootstrap.Config) bootstrap.Db,
) func() (master, slave *redis.Client, resultErr error) {
	return func() (master *redis.Client, slave *redis.Client, resultErr error) {
		cfg, err := configGetterFn()
		if err != nil {
			resultErr = err
			return
		}

		dbCfg := configSelector(cfg)
		master, resultErr = conn.Connect(dbCfg)
		if resultErr != nil {
			return
		}
		setConnect(cfg.RedisConfig, master)

		slave, resultErr = conn.Connect(dbCfg)
		if resultErr != nil {
			return
		}
		setConnect(cfg.RedisConfig, slave)
		return
	}
}

func setConnect(connCfg bootstrap.RedisConfig, connection *redis.Client) {
	maxLifeHour := connCfg.MaxLifeHour
	maxConnAge := time.Hour * time.Duration(maxLifeHour)

	connection.Options().MaxConnAge = maxConnAge
}

func Dispose() {
	if tradeDb != nil {
		tradeDb.Dispose()
	}
}
