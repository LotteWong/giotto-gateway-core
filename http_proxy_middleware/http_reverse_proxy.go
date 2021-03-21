package http_proxy_middleware

import (
	"github.com/LotteWong/giotto-gateway/middleware"
	"github.com/LotteWong/giotto-gateway/models/po"
	"github.com/LotteWong/giotto-gateway/reverse_proxy"
	"github.com/LotteWong/giotto-gateway/service"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

func HttpReverseProxyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		httpServiceInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, http.StatusInternalServerError, errors.New("service not found"))
			c.Abort()
			return
		}
		httpServiceDetail := httpServiceInterface.(*po.ServiceDetail)

		lb, err := service.GetLbService().GetLbWithConfForSvc(httpServiceDetail)
		if err != nil {
			middleware.ResponseError(c, http.StatusInternalServerError, err)
			c.Abort()
			return
		}
		trans, err := service.GetTransService().GetTransForSvc(httpServiceDetail)
		if err != nil {
			middleware.ResponseError(c, http.StatusInternalServerError, err)
			c.Abort()
			return
		}

		// use reverse proxy to serve http
		proxy := reverse_proxy.NewReverseProxy(c, lb, trans)
		proxy.ServeHTTP(c.Writer, c.Request)

		// abort the original server to be accessed
		c.Abort()
		return
	}
}
