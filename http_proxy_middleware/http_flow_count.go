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

func HttpFlowCountMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		httpServiceInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, http.StatusInternalServerError, errors.New("service not found"))
			c.Abort()
			return
		}
		httpServiceDetail := httpServiceInterface.(*po.ServiceDetail)

		appInterface, ok := c.Get("app")
		if !ok {
			middleware.ResponseError(c, http.StatusInternalServerError, errors.New("app not found"))
			c.Abort()
			return
		}
		app := appInterface.(*po.App)

		ttlServiceName := constants.TotalFlowCountPrefix
		ttlFlowCount, err := service.GetFlowCountService().GetFlowCount(ttlServiceName)
		if err != nil {
			middleware.ResponseError(c, http.StatusInternalServerError, err)
			c.Abort()
			return
		}
		service.GetFlowCountService().Increase(ttlFlowCount)

		svcServiceName := constants.ServiceFlowCountPrefix + httpServiceDetail.Info.ServiceName
		svcFlowCount, err := service.GetFlowCountService().GetFlowCount(svcServiceName)
		if err != nil {
			middleware.ResponseError(c, http.StatusInternalServerError, err)
			c.Abort()
			return
		}
		service.GetFlowCountService().Increase(svcFlowCount)

		appServiceName := constants.AppFlowCountPrefix + app.AppId
		appFlowCount, err := service.GetFlowCountService().GetFlowCount(appServiceName)
		if err != nil {
			middleware.ResponseError(c, http.StatusInternalServerError, err)
			c.Abort()
			return
		}
		service.GetFlowCountService().Increase(appFlowCount)
		if app.Qpd > 0 && appFlowCount.TotalCount > app.Qpd {
			middleware.ResponseError(c, http.StatusInternalServerError, errors.New(fmt.Sprintf("app's qpd exceeds limit, current: %d, limit: %d", appFlowCount.TotalCount, app.Qpd)))
			c.Abort()
			return
		}

		c.Next()
	}
}
