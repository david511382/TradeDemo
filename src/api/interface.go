package api

import (
	"zerologix-homework/src/server/reqs"
	"zerologix-homework/src/server/resp"
)

type IOrder interface {
	PostBuy(req *reqs.OrderPostBuy) (
		result resp.OrderPostBuy,
		resultErr error,
	)
	PostSell(req *reqs.OrderPostSell) (
		result resp.OrderPostSell,
		resultErr error,
	)
}

type IMatchOrder interface {
	Match(price float64) (
		resultErr error,
	)
}
