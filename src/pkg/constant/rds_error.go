package constant

import (
	errUtil "zerologix-homework/src/pkg/util/error"

	"github.com/rs/zerolog"
)

const (
	ERROR_MSG_REDIS_NOT_CHANGE  = "redis Not Change"
	ERROR_MSG_REDIS_NOT_SUCCESS = "redis Not Success"
	ERROR_MSG_REDIS_NO_DATA     = "redis No Data"
	ERROR_MSG_REDIS_NOT_EXIST   = "redis: nil"
)

var (
	ErrInfoRedisNotChange = errUtil.New(ERROR_MSG_REDIS_NOT_CHANGE, zerolog.InfoLevel)
)
