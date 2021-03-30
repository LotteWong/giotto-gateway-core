package http_proxy_middleware

import (
	"github.com/LotteWong/giotto-gateway/common_middleware"
	"github.com/LotteWong/giotto-gateway/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HttpProxyAccessMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		httpService, err := service.GetSvcService().HttpProxyAccessService(c)
		if err != nil {
			common_middleware.ResponseError(c, http.StatusInternalServerError, err)
			c.Abort()
			return
		}
		//log.Printf("matched http service: %s", utils.Obj2Json(httpService))
		c.Set("service", httpService)
		c.Next()
	}
}
