package http_proxy_middleware

import (
	"net/http"

	"github.com/LotteWong/giotto-gateway-core/common_middleware"
	"github.com/LotteWong/giotto-gateway-core/constants"
	"github.com/LotteWong/giotto-gateway-core/models/po"
	"github.com/LotteWong/giotto-gateway-core/reverse_proxy"
	"github.com/LotteWong/giotto-gateway-core/service"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func HttpReverseProxyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		httpServiceInterface, ok := c.Get("service")
		if !ok {
			common_middleware.ResponseError(c, http.StatusInternalServerError, errors.New("service not found"))
			c.Abort()
			return
		}
		httpServiceDetail := httpServiceInterface.(*po.ServiceDetail)

		lb, err := service.GetLbService().GetLbWithConfForSvc(httpServiceDetail)
		if err != nil {
			common_middleware.ResponseError(c, http.StatusInternalServerError, err)
			c.Abort()
			return
		}
		trans, err := service.GetTransService().GetTransForSvc(httpServiceDetail)
		if err != nil {
			common_middleware.ResponseError(c, http.StatusInternalServerError, err)
			c.Abort()
			return
		}

		// use reverse proxy to serve http
		var scheme string
		if httpServiceDetail.HttpRule.NeedHttps == constants.Enable {
			scheme = "https://"
		} else {
			scheme = "http://"
		}
		proxy := reverse_proxy.NewHttpReverseProxy(c, lb, trans, scheme)
		proxy.ServeHTTP(c.Writer, c.Request)

		// abort the original server to be accessed
		c.Abort()
		return
	}
}
