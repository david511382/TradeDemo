package trade

type ITradeLogic interface {
	Match(order Order) (
		resultErr error,
	)
}

type ITradeStorage interface {
	LockID() (
		resultErr error,
	)
	ReleaseIDLock()
	LoadID() (
		id uint,
		resultErr error,
	)
	UpdateID(
		id uint,
	) (
		resultErr error,
	)

	LockOrders(
		price float64,
	) (
		resultErr error,
	)
	ReleaseOrdersLock(
		price float64,
	)
	Load(
		price float64,
	) (
		storeOrders []*Order,
		resultErr error,
	)
	Set(
		price float64,
		orders []*Order,
	) (
		resultErr error,
	)
	Delete(
		price float64,
	) (
		resultErr error,
	)
}
