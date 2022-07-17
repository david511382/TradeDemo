package trade

type OrderType uint8

const (
	ORDER_TYPE_BUY  OrderType = 1
	ORDER_TYPE_SELL OrderType = 2
)

type Order struct {
	ID        uint
	OrderType OrderType
	Quantity  uint
	Price     float64
	Timestamp int64
}
