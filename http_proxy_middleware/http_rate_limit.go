package http_proxy_middleware

import (
	"fmt"
	"github.com/LotteWong/giotto-gateway/constants"
	"github.com/LotteWong/giotto-gateway/middleware"
	"github.com/LotteWong/giotto-gateway/models/po"
	"github.com/LotteWong/giotto-gateway/service"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

func HttpRateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		httpServiceInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, http.StatusInternalServerError, errors.New("service not found"))
			c.Abort()
			return
		}
		httpServiceDetail := httpServiceInterface.(*po.ServiceDetail)

		if httpServiceDetail.AccessControl.ServiceHostFlowLimit != 0 {
			svrServiceName := constants.ServiceFlowCountPrefix + httpServiceDetail.Info.ServiceName
			svrRateLimit, err := service.GetRateLimitService().GetRateLimit(svrServiceName, httpServiceDetail.AccessControl.ServiceHostFlowLimit)
			if err != nil {
				middleware.ResponseError(c, http.StatusInternalServerError, err)
				c.Abort()
				return
			}
			if !svrRateLimit.Allow() {
				middleware.ResponseError(c, http.StatusInternalServerError, errors.New(fmt.Sprintf("service host flow limit is %d, rate limit exceeds", httpServiceDetail.AccessControl.ServiceHostFlowLimit)))
				c.Abort()
				return
			}
		}

		if httpServiceDetail.AccessControl.ClientIpFlowLimit != 0 {
			cltServiceName := constants.ServiceFlowCountPrefix + httpServiceDetail.Info.ServiceName + "_" + c.ClientIP()
			cltRateLimit, err := service.GetRateLimitService().GetRateLimit(cltServiceName, httpServiceDetail.AccessControl.ClientIpFlowLimit)
			if err != nil {
				middleware.ResponseError(c, http.StatusInternalServerError, err)
				c.Abort()
				return
			}
			if !cltRateLimit.Allow() {
				middleware.ResponseError(c, http.StatusInternalServerError, errors.New(fmt.Sprintf("client ip flow limit is %d, rate limit exceeds", httpServiceDetail.AccessControl.ClientIpFlowLimit)))
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
