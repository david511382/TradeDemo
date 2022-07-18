package router

import (
	"zerologix-homework/bootstrap"
	"zerologix-homework/src/server/api/trade/order"
	"zerologix-homework/src/server/domain"

	"github.com/gin-gonic/gin"
)

func SetupCustomerRouter(cfg *bootstrap.Config, tokenVerifier domain.ITokenVerifier, api *gin.RouterGroup) {
	apiTrade := api.Group("/trade")

	apiTradeOrder := apiTrade.Group("/order")
	// api/trade/order/test
	apiTradeOrder.POST("/test", order.PostTest)
	// api/trade/order/buy
	apiTradeOrder.POST("/buy", order.PostBuy)
	// api/trade/order/sell
	apiTradeOrder.POST("/sell", order.PostSell)
}
