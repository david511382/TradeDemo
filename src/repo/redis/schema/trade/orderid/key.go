package orderid

import (
	"strconv"
	"zerologix-homework/src/repo/redis/common"
)

type Key struct {
	common.BaseKey[uint]
}

func New(connectionCreator common.IConnection, baseKey string) *Key {
	result := &Key{}
	result.BaseKey = *common.NewBaseKey[uint](
		connectionCreator,
		baseKey+"LastOrderID",
		result,
	)
	return result
}

func (k Key) StringifyValue(value uint) (string, error) {
	return strconv.FormatUint(uint64(value), 10), nil
}

func (k Key) ParseValue(valueStr string) (uint, error) {
	u64, err := strconv.ParseUint(valueStr, 10, 64)
	return uint(u64), err
}
