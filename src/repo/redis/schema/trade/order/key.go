package order

import (
	"encoding/json"
	"strconv"
	"zerologix-homework/src/repo/redis/common"
)

type Key struct {
	common.BaseHashKey[float64, []*Model]
}

func New(connectionCreator common.IConnection, baseKey string) *Key {
	result := &Key{}
	result.BaseHashKey = *common.NewBaseHashKey[float64, []*Model](
		connectionCreator,
		baseKey+"Order",
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

func (k Key) StringifyValue(value []*Model) (string, error) {
	bs, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func (k Key) ParseValue(valueStr string) ([]*Model, error) {
	value := make([]*Model, 0)
	err := json.Unmarshal([]byte(valueStr), &value)
	if err != nil {
		return nil, err
	}
	return value, nil
}

type Model struct {
	ID        uint
	OrderType uint8
	Quantity  uint
	Timestamp int64
}
