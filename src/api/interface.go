package api

import (
	"zerologix-homework/src/server/reqs"
	"zerologix-homework/src/server/resp"
)

type IOrder interface {
	PostTest(req *reqs.OrderPostTest) (
		result resp.OrderPostTest,
		resultErr error,
	)
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
