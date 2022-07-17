package orderidlocker

import (
	"strconv"
	"zerologix-homework/src/repo/redis/common"
)

type Key struct {
	common.BaseKey[int64]
}

func New(connectionCreator common.IConnection, baseKey string) *Key {
	result := &Key{}
	result.BaseKey = *common.NewBaseKey[int64](
		connectionCreator,
		baseKey+"LastOrderIDLocker",
		result,
	)
	return result
}

func (k Key) StringifyValue(value int64) (string, error) {
	return strconv.FormatInt(value, 10), nil
}

func (k Key) ParseValue(valueStr string) (int64, error) {
	return strconv.ParseInt(valueStr, 10, 64)
}
