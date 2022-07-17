package conn

import (
	"zerologix-homework/bootstrap"
	errUtil "zerologix-homework/src/pkg/util/error"

	"github.com/go-redis/redis"
)

func Connect(cfg bootstrap.Db) (*redis.Client, error) {
	url := cfg.ParseToUrl()
	rdsOpt, err := redis.ParseURL(url)
	if err != nil {
		return nil, errUtil.NewError(err)
	}
	connection := redis.NewClient(rdsOpt)

	if err := connection.Ping().Err(); err != nil {
		return nil, errUtil.NewError(err)
	}

	return connection, nil
}
