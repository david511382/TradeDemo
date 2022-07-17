package trade

import (
	"zerologix-homework/src/repo/redis/common"
	"zerologix-homework/src/repo/redis/schema/trade/order"
	"zerologix-homework/src/repo/redis/schema/trade/orderid"
	"zerologix-homework/src/repo/redis/schema/trade/orderidlocker"
	"zerologix-homework/src/repo/redis/schema/trade/orderlocker"

	"github.com/go-redis/redis"
)

type Database struct {
	*common.BaseDatabase[Schema]

	Schema
}

func NewDatabase(connect func() (master, slave *redis.Client, resultErr error), baseKey string) *Database {
	result := &Database{
		BaseDatabase: common.NewBaseDatabase(
			connect,
			func(connectionCreator common.IConnection) Schema {
				return *NewSchema(connectionCreator, baseKey)
			},
			baseKey,
		),
	}
	result.Schema = *NewSchema(result, baseKey)
	return result
}

func SchemaCreator(connectionCreator common.IConnection, baseKey string) Schema {
	return *NewSchema(connectionCreator, baseKey)
}

type Schema struct {
	Order         *order.Key
	OrderLocker   *orderlocker.Key
	OrderID       *orderid.Key
	OrderIDLocker *orderidlocker.Key
}

func NewSchema(connectionCreator common.IConnection, baseKey string) *Schema {
	result := &Schema{
		Order:         order.New(connectionCreator, baseKey),
		OrderLocker:   orderlocker.New(connectionCreator, baseKey),
		OrderID:       orderid.New(connectionCreator, baseKey),
		OrderIDLocker: orderidlocker.New(connectionCreator, baseKey),
	}
	return result
}
