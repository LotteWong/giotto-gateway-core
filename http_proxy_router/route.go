package router

import (
	"github.com/LotteWong/giotto-gateway-core/common_middleware"
	"github.com/LotteWong/giotto-gateway-core/controller"
	"github.com/LotteWong/giotto-gateway-core/http_proxy_middleware"
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

	// jwt api routes
	jwtGroup := router.Group("/tokens")
	jwtGroup.Use(common_middleware.TranslationMiddleware())
	{
		// Post    /tokens/jwt
		controller.RegistJwtRoutes(jwtGroup)
	}

	router.Use(
		http_proxy_middleware.HttpProxyAccessMiddleware(),

		http_proxy_middleware.HttpJwtAuthMiddleware(),

		http_proxy_middleware.HttpFlowCountMiddleware(),
		http_proxy_middleware.HttpRateLimitMiddleware(),

		http_proxy_middleware.HttpIpAuthMiddleware(),

		http_proxy_middleware.HttpRouteRewriteMiddleware(),

		http_proxy_middleware.HttpReverseProxyMiddleware(),
	)

	return router
}
