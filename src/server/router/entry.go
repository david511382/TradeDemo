package router

import (
	"io"
	"zerologix-homework/bootstrap"
	"zerologix-homework/src/logger"
	"zerologix-homework/src/server/common"
	"zerologix-homework/src/server/middleware"
	"zerologix-homework/src/server/swagger"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func EntryRouter(cfg *bootstrap.Config) *gin.Engine {
	if cfg.Var.UseDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 取消打印文字顏色
	gin.DisableConsoleColor()
	// 使用打印文字顏色
	gin.ForceConsoleColor()

	// 設定gin
	router := NewRouter()

	router.Use(middleware.Logger())
	router.Use(middleware.Cors)

	// docs
	doc := router.Group("/docs")
	doc.GET("/*any", swagger.Swag)

	router.Use(gin.Logger())

	// api
	SetupApiRouter(cfg, router)

	return router
}

func NewRouter() *gin.Engine {
	// 設定gin
	router := gin.New()

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.RecoveryWithWriter(io.MultiWriter(
		logger.GetWriter(logger.NAME_API, zerolog.ErrorLevel),
	)))

	return router
}

func SetupApiRouter(cfg *bootstrap.Config, router *gin.Engine) *gin.Engine {
	defaultTokenVerifier := common.NewDebugTokenVerifier()

	// api
	api := router.Group("/api")

	SetupCustomerRouter(cfg, defaultTokenVerifier, api)

	return router
}
