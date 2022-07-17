package order

import (
	"zerologix-homework/src/api"
	tradeBll "zerologix-homework/src/bll/trade"
	"zerologix-homework/src/pkg/timeutil"
	errUtil "zerologix-homework/src/pkg/util/error"
	"zerologix-homework/src/repo/redis"
	"zerologix-homework/src/server/common"
	"zerologix-homework/src/server/reqs"

	"github.com/gin-gonic/gin"
)

// PostBuy 買
// @Tags Order
// @Summary 買
// @Description 買
// @Accept json
// @Produce json
// @Param param body reqs.OrderPostBuy true "參數"
// @Success 200 {object} resp.Base{} "資料"
// @Security ApiKeyAuth
// @Router /trade/order/buy [post]
func PostBuy(c *gin.Context) {
	tradeRds := redis.Trade()
	timeUtil, err := timeutil.GetTimeUtil()
	if err != nil {
		common.Fail(c, err)
		return
	}
	tradeStorage := tradeBll.NewTradeStorage(
		timeUtil,
		tradeRds,
	)
	tradeLogic := tradeBll.NewTradeLogic(
		timeUtil,
		tradeStorage,
	)
	badmintonActivityApiLogic := api.NewOrder(
		tradeLogic,
	)
	NewPostBuyHandler(badmintonActivityApiLogic)(c)
}

func NewPostBuyHandler(logic api.IOrder) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			req reqs.OrderPostBuy
		)
		if err := c.ShouldBindJSON(&req); err != nil {
			err := errUtil.NewError(err)
			common.FailRequest(c, err)
			return
		}

		result, err := logic.PostBuy(&req)
		if err != nil {
			common.Fail(c, err)
			return
		}

		common.Success(c, result)
	}
}

// PostSell 賣
// @Tags Order
// @Summary 賣
// @Description 賣
// @Accept json
// @Produce json
// @Param param body reqs.OrderPostSell true "參數"
// @Success 200 {object} resp.Base{} "資料"
// @Security ApiKeyAuth
// @Router /trade/order/sell [post]
func PostSell(c *gin.Context) {
	tradeRds := redis.Trade()
	timeUtil, err := timeutil.GetTimeUtil()
	if err != nil {
		common.Fail(c, err)
		return
	}
	tradeStorage := tradeBll.NewTradeStorage(
		timeUtil,
		tradeRds,
	)
	tradeLogic := tradeBll.NewTradeLogic(
		timeUtil,
		tradeStorage,
	)
	badmintonActivityApiLogic := api.NewOrder(
		tradeLogic,
	)
	NewPostSellHandler(badmintonActivityApiLogic)(c)
}

func NewPostSellHandler(logic api.IOrder) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			req reqs.OrderPostSell
		)
		if err := c.ShouldBindJSON(&req); err != nil {
			err := errUtil.NewError(err)
			common.FailRequest(c, err)
			return
		}

		result, err := logic.PostSell(&req)
		if err != nil {
			common.Fail(c, err)
			return
		}

		common.Success(c, result)
	}
}
