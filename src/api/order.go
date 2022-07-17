package api

import (
	tradeBll "zerologix-homework/src/bll/trade"
	"zerologix-homework/src/server/reqs"
	"zerologix-homework/src/server/resp"
)

type Order struct {
	tradeLogic tradeBll.ITradeLogic
}

func NewOrder(
	tradeLogic tradeBll.ITradeLogic,
) *Order {
	return &Order{
		tradeLogic: tradeLogic,
	}
}

func (l Order) PostBuy(req *reqs.OrderPostBuy) (
	result resp.OrderPostBuy,
	resultErr error,
) {
	if err := l.tradeLogic.Match(tradeBll.Order{
		OrderType: tradeBll.ORDER_TYPE_BUY,
		Quantity:  req.Quantity,
		Price:     req.Price,
	}); err != nil {
		resultErr = err
		result.Message = "Error"
		return
	}

	result.Message = "OK"

	return
}

func (l Order) PostSell(req *reqs.OrderPostSell) (
	result resp.OrderPostSell,
	resultErr error,
) {
	if err := l.tradeLogic.Match(tradeBll.Order{
		OrderType: tradeBll.ORDER_TYPE_SELL,
		Quantity:  req.Quantity,
		Price:     req.Price,
	}); err != nil {
		resultErr = err
		result.Message = "Error"
		return
	}

	result.Message = "OK"

	return
}
