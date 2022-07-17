package orderlocker

import (
	"strconv"
	"zerologix-homework/src/repo/redis/common"
)

type Key struct {
	common.BaseHashKey[float64, int64]
}

func New(connectionCreator common.IConnection, baseKey string) *Key {
	result := &Key{}
	result.BaseHashKey = *common.NewBaseHashKey[float64, int64](
		connectionCreator,
		baseKey+"OrderLocker",
		result,
	)
	return result
}

func (k Key) StringifyField(price float64) string {
	return strconv.FormatFloat(price, 'e', 64, 64)
}

func (k Key) ParseField(fieldStr string) (price float64, resultErr error) {
	return strconv.ParseFloat(fieldStr, 64)
}

func (k Key) StringifyValue(value int64) (string, error) {
	return strconv.FormatInt(value, 10), nil
}

func (k Key) ParseValue(valueStr string) (int64, error) {
	return strconv.ParseInt(valueStr, 10, 64)
}
