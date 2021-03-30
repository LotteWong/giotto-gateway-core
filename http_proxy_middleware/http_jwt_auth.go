package http_proxy_middleware

import (
	"fmt"
	"github.com/LotteWong/giotto-gateway/common_middleware"
	"github.com/LotteWong/giotto-gateway/constants"
	"github.com/LotteWong/giotto-gateway/models/po"
	"github.com/LotteWong/giotto-gateway/service"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

func HttpJwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		httpServiceInterface, ok := c.Get("service")
		if !ok {
			common_middleware.ResponseError(c, http.StatusInternalServerError, errors.New("service not found"))
			c.Abort()
			return
		}
		httpServiceDetail := httpServiceInterface.(*po.ServiceDetail)

		// parse authorization to get jwt
		info := strings.Split(c.GetHeader("Authorization"), " ")
		if len(info) != 2 {
			common_middleware.ResponseError(c, http.StatusInternalServerError, errors.New("can not get jwt from authorization header"))
			c.Abort()
			return
		}

		tokenType := info[0]
		tokenString := info[1]

		// verify jwt by expire at and issuer
		var err error
		switch tokenType {
		case constants.JwtType:
			err = service.GetJwtService().VerifyJwt(c, httpServiceDetail, tokenString)
		default:
			err = errors.New(fmt.Sprintf("not support jwt type %s", tokenType))
		}

		if err != nil {
			common_middleware.ResponseError(c, http.StatusInternalServerError, err)
			c.Abort()
			return
		}
		c.Next()
	}
}
