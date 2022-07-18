package trade

import (
	"sort"
	"zerologix-homework/src/pkg/timeutil"
)

type TradeLogic struct {
	timeUtil     timeutil.ITime
	tradeStorage ITradeStorage
}

func NewTradeLogic(
	timeUtil timeutil.ITime,
	tradeStorage ITradeStorage,
) *TradeLogic {
	return &TradeLogic{
		timeUtil:     timeUtil,
		tradeStorage: tradeStorage,
	}
}

func (l *TradeLogic) Match(order Order) (
	resultErr error,
) {
	price := order.Price

	// lock
	{
		err := l.tradeStorage.LockOrders(price)
		defer func() {
			l.tradeStorage.ReleaseOrdersLock(price)
		}()
		if err != nil {
			resultErr = err
			return
		}
	}

	storeOrders, err := l.tradeStorage.Load(
		price,
	)
	if err != nil {
		resultErr = err
		return
	}

	sort.Slice(storeOrders, func(i, j int) bool {
		return storeOrders[i].Timestamp < storeOrders[j].Timestamp
	})

	var matchOrderType OrderType
	switch order.OrderType {
	case ORDER_TYPE_BUY:
		matchOrderType = ORDER_TYPE_SELL
	case ORDER_TYPE_SELL:
		matchOrderType = ORDER_TYPE_BUY
	}

	resultStoreOrders, isDone := l.matchOrder(
		matchOrderType,
		order.Quantity,
		storeOrders,
	)
	if !isDone {
		newID, err := l.getNewID()
		if err != nil {
			resultErr = err
			return
		}
		timestamp := l.timeUtil.Now().UnixNano()

		order.ID = newID
		order.Timestamp = timestamp
		resultStoreOrders = append(resultStoreOrders, &order)
	}

	if len(resultStoreOrders) == 0 {
		if err := l.tradeStorage.Delete(price); err != nil {
			resultErr = err
			return
		}
	} else {
		if err := l.tradeStorage.Set(price, resultStoreOrders); err != nil {
			resultErr = err
			return
		}
	}

	return
}

func (l *TradeLogic) getNewID() (uint, error) {
	// lock
	{
		err := l.tradeStorage.LockID()
		defer func() {
			l.tradeStorage.ReleaseIDLock()
		}()
		if err != nil {
			return 0, err
		}
	}

	currentID, err := l.tradeStorage.LoadID()
	if err != nil {
		return 0, err
	}

	newID := currentID + 1

	if err := l.tradeStorage.UpdateID(newID); err != nil {
		return 0, err
	}

	return newID, nil
}

func (l *TradeLogic) matchOrder(
	matchOrderType OrderType,
	matchQuantity uint,
	matchOrders []*Order,
) (
	resultOrders []*Order,
	isDone bool,
) {
	for i, v := range matchOrders {
		if matchOrderType != v.OrderType {
			resultOrders = append(resultOrders, v)
			continue
		}

		if v.Quantity > matchQuantity {
			leftQuantity := v.Quantity - matchQuantity
			v.Quantity = leftQuantity
			matchQuantity = 0
			resultOrders = append(resultOrders, matchOrders[i:]...)
			break
		}

		useQuantity := v.Quantity
		matchQuantity -= useQuantity
		v.Quantity = useQuantity

		if matchQuantity == 0 {
			nextIndex := i + 1
			if nextIndex < len(matchOrders) {
				resultOrders = append(resultOrders, matchOrders[nextIndex:]...)
			}
			break
		}
	}

	isDone = matchQuantity == 0

	return
}
