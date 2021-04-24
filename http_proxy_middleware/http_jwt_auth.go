package http_proxy_middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/LotteWong/giotto-gateway-core/common_middleware"
	"github.com/LotteWong/giotto-gateway-core/constants"
	"github.com/LotteWong/giotto-gateway-core/models/po"
	"github.com/LotteWong/giotto-gateway-core/service"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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
		pair := strings.Split(c.GetHeader("Authorization"), " ")
		if len(pair) != 2 {
			common_middleware.ResponseError(c, http.StatusInternalServerError, errors.New("can not get jwt from authorization header"))
			c.Abort()
			return
		}

		tokenType := pair[0]
		tokenString := pair[1]

		// verify jwt by expire at and issuer
		var err error
		switch tokenType {
		case constants.JwtType:
			err = service.GetJwtService().HttpVerifyJwt(c, httpServiceDetail, tokenString)
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
