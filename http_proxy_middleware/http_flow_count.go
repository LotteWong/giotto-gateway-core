package http_proxy_middleware

import (
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

		c.Next()
	}
}
