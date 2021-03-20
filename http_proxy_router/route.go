package router

import (
	"github.com/LotteWong/giotto-gateway/http_proxy_middleware"
	"github.com/LotteWong/giotto-gateway/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares...)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	commonMiddlewares := []gin.HandlerFunc{
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
		http_proxy_middleware.HttpProxyAccessMiddleware(),
	}
	router.Use(commonMiddlewares...)

	return router
}
