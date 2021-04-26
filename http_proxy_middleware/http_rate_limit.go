package http_proxy_middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/LotteWong/giotto-gateway-core/common_middleware"
	"github.com/LotteWong/giotto-gateway-core/constants"
	"github.com/LotteWong/giotto-gateway-core/models/po"
	"github.com/LotteWong/giotto-gateway-core/service"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func HttpRateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		httpServiceInterface, ok := c.Get("service")
		if !ok {
			common_middleware.ResponseError(c, http.StatusInternalServerError, errors.New("service not found"))
			c.Abort()
			return
		}
		httpServiceDetail := httpServiceInterface.(*po.ServiceDetail)

		appInterface, ok := c.Get("app")
		if !ok {
			common_middleware.ResponseError(c, http.StatusInternalServerError, errors.New("app not found"))
			c.Abort()
			return
		}
		app := appInterface.(*po.App)

		if httpServiceDetail.AccessControl.ServiceHostFlowLimit != 0 {
			svrServiceName := constants.ServiceFlowCountPrefix + httpServiceDetail.Info.ServiceName
			svrRateLimit, err := service.GetRateLimitService().GetRateLimit(svrServiceName)
			if err != nil {
				common_middleware.ResponseError(c, http.StatusInternalServerError, err)
				c.Abort()
				return
			}
			_, _, svrAllow := svrRateLimit.Allow(svrServiceName, httpServiceDetail.AccessControl.ServiceHostFlowLimit, 1*time.Second)
			// log.Printf("svr name:%s, count:%d\n", svrServiceName, svrCount)
			if !svrAllow {
				common_middleware.ResponseError(c, http.StatusInternalServerError, errors.New(fmt.Sprintf("service host flow limit is %d, rate limit exceeds", httpServiceDetail.AccessControl.ServiceHostFlowLimit)))
				c.Abort()
				return
			}
		}

		if httpServiceDetail.AccessControl.ClientIpFlowLimit != 0 {
			cltServiceName := constants.ServiceFlowCountPrefix + httpServiceDetail.Info.ServiceName + "_" + c.ClientIP()
			cltRateLimit, err := service.GetRateLimitService().GetRateLimit(cltServiceName)
			if err != nil {
				common_middleware.ResponseError(c, http.StatusInternalServerError, err)
				c.Abort()
				return
			}
			_, _, cltAllow := cltRateLimit.Allow(cltServiceName, httpServiceDetail.AccessControl.ClientIpFlowLimit, 1*time.Second)
			// log.Printf("clt name:%s, count:%d\n", cltServiceName, cltCount)
			if !cltAllow {
				common_middleware.ResponseError(c, http.StatusInternalServerError, errors.New(fmt.Sprintf("client ip flow limit is %d, rate limit exceeds", httpServiceDetail.AccessControl.ClientIpFlowLimit)))
				c.Abort()
				return
			}
		}

		if app.Qps != 0 {
			appServiceName := constants.AppFlowCountPrefix + app.AppId
			appRateLimit, err := service.GetRateLimitService().GetRateLimit(appServiceName)
			if err != nil {
				common_middleware.ResponseError(c, http.StatusInternalServerError, err)
				c.Abort()
				return
			}
			_, _, appAllow := appRateLimit.Allow(appServiceName, app.Qps, 1*time.Second)
			// log.Printf("app name:%s, count:%d\n", appServiceName, appCount)
			if !appAllow {
				common_middleware.ResponseError(c, http.StatusInternalServerError, errors.New(fmt.Sprintf("app flow limit is %d, rate limit exceeds", app.Qps)))
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
