package http_proxy_middleware

import (
	"fmt"
	"github.com/LotteWong/giotto-gateway/common_middleware"
	"github.com/LotteWong/giotto-gateway/constants"
	"github.com/LotteWong/giotto-gateway/models/po"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

func HttpIpAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		httpServiceInterface, ok := c.Get("service")
		if !ok {
			common_middleware.ResponseError(c, http.StatusInternalServerError, errors.New("service not found"))
			c.Abort()
			return
		}
		httpServiceDetail := httpServiceInterface.(*po.ServiceDetail)

		var whiteIpList []string
		var blackIpList []string
		openAuth := httpServiceDetail.AccessControl.OpenAuth
		whiteList := httpServiceDetail.AccessControl.WhiteList
		if whiteList != "" {
			whiteIpList = strings.Split(whiteList, ",")
		}
		blackList := httpServiceDetail.AccessControl.BlackList
		if blackList != "" {
			blackIpList = strings.Split(blackList, ",")
		}

		if openAuth == constants.Enable {
			if len(whiteIpList) > 0 { // white list has higher priority
				if !checkStrInSlice(whiteIpList, c.ClientIP()) {
					common_middleware.ResponseError(c, http.StatusInternalServerError, errors.New(fmt.Sprintf("ip %s not in white ip list", c.ClientIP())))
					c.Abort()
					return
				}
			} else { // black list has lower priority
				if len(blackIpList) > 0 {
					if checkStrInSlice(blackIpList, c.ClientIP()) {
						common_middleware.ResponseError(c, http.StatusInternalServerError, errors.New(fmt.Sprintf("ip %s is in black ip list", c.ClientIP())))
						c.Abort()
						return
					}
				}
			}
		}

		c.Next()
	}
}

func checkStrInSlice(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}
