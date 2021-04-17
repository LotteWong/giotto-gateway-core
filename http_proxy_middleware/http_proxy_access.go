package http_proxy_middleware

import (
	"log"
	"net/http"

	"github.com/LotteWong/giotto-gateway/common_middleware"
	"github.com/LotteWong/giotto-gateway/service"
	"github.com/LotteWong/giotto-gateway/utils"
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
		log.Printf("matched http service: %s\n", utils.Obj2Json(httpService))
		c.Set("service", httpService)
		c.Next()
	}
}
