package tcp_proxy_middleware

import (
	"fmt"
	"strings"

	"github.com/LotteWong/giotto-gateway/constants"
	"github.com/LotteWong/giotto-gateway/models/po"
	tcp_server "github.com/LotteWong/tcp-conn-server"
	"github.com/pkg/errors"
)

func TcpIpAuthMiddleware() func(c *tcp_server.TcpRouterContext) {
	return func(c *tcp_server.TcpRouterContext) {
		tcpServiceInterface := c.Get("service")
		if tcpServiceInterface == nil {
			c.Conn.Write([]byte("service not found"))
			c.Abort()
			return
		}
		tcpServiceDetail := tcpServiceInterface.(*po.ServiceDetail)

		var whiteIpList []string
		var blackIpList []string
		openAuth := tcpServiceDetail.AccessControl.OpenAuth
		whiteList := tcpServiceDetail.AccessControl.WhiteList
		if whiteList != "" {
			whiteIpList = strings.Split(whiteList, ",")
		}
		blackList := tcpServiceDetail.AccessControl.BlackList
		if blackList != "" {
			blackIpList = strings.Split(blackList, ",")
		}

		if openAuth == constants.Enable {
			var clientIp string
			pair := strings.Split(c.Conn.RemoteAddr().String(), ":")
			if len(pair) != 2 {
				c.Conn.Write([]byte("can not get client ip and port"))
				c.Abort()
				return
			}
			clientIp = pair[0]

			if len(whiteIpList) > 0 { // white list has higher priority
				if !checkStrInSlice(whiteIpList, clientIp) {
					c.Conn.Write([]byte(errors.New(fmt.Sprintf("ip %s not in white ip list", clientIp)).Error()))
					c.Abort()
					return
				}
			} else { // black list has lower priority
				if len(blackIpList) > 0 {
					if checkStrInSlice(blackIpList, clientIp) {
						c.Conn.Write([]byte(errors.New(fmt.Sprintf("ip %s is in black ip list", clientIp)).Error()))
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
