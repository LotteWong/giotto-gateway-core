package http_proxy_middleware

import (
	"net/http"

	"github.com/LotteWong/giotto-gateway-core/common_middleware"
	"github.com/LotteWong/giotto-gateway-core/service"
	"github.com/gin-gonic/gin"
)

func HttpProxyAccessMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		httpService, err := service.GetSvcService().HttpProxyAccessService(c)
		if err != nil {
			common_middleware.ResponseError(c, http.StatusInternalServerError, err)
			c.Abort()
			return
		}
		c.Set("service", httpService)

		scheme, err := service.GetSvcService().HttpProxyAccessScheme(c, httpService)
		if err != nil {
			common_middleware.ResponseError(c, http.StatusInternalServerError, err)
			c.Abort()
			return
		}
		c.Set("scheme", scheme)

		c.Next()
	}
}
