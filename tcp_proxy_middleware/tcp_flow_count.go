package tcp_proxy_middleware

import (
	"github.com/LotteWong/giotto-gateway/constants"
	"github.com/LotteWong/giotto-gateway/models/po"
	"github.com/LotteWong/giotto-gateway/service"
	tcp "github.com/LotteWong/tcp-conn-server"
)

func TcpFlowCountMiddleware() func(c *tcp.TcpRouterContext) {
	return func(c *tcp.TcpRouterContext) {
		tcpServiceInterface := c.Get("service")
		if tcpServiceInterface == nil {
			c.Conn.Write([]byte("service not found"))
			c.Abort()
			return
		}
		tcpServiceDetail := tcpServiceInterface.(*po.ServiceDetail)

		// appInterface := c.Get("app")
		// if appInterface == nil {
		// 	c.Conn.Write([]byte("app not found"))
		// 	c.Abort()
		// 	return
		// }
		// app := appInterface.(*po.App)

		ttlServiceName := constants.TotalFlowCountPrefix
		ttlFlowCount, err := service.GetFlowCountService().GetFlowCount(ttlServiceName)
		if err != nil {
			c.Conn.Write([]byte(err.Error()))
			c.Abort()
			return
		}
		service.GetFlowCountService().Increase(ttlFlowCount)

		svcServiceName := constants.ServiceFlowCountPrefix + tcpServiceDetail.Info.ServiceName
		svcFlowCount, err := service.GetFlowCountService().GetFlowCount(svcServiceName)
		if err != nil {
			c.Conn.Write([]byte(err.Error()))
			c.Abort()
			return
		}
		service.GetFlowCountService().Increase(svcFlowCount)

		// appServiceName := constants.AppFlowCountPrefix + app.AppId
		// appFlowCount, err := service.GetFlowCountService().GetFlowCount(appServiceName)
		// if err != nil {
		// 	c.Conn.Write([]byte(err.Error()))
		// 	c.Abort()
		// 	return
		// }
		// service.GetFlowCountService().Increase(appFlowCount)
		// if app.Qpd > 0 && appFlowCount.TotalCount > app.Qpd {
		// 	c.Conn.Write([]byte(errors.New(fmt.Sprintf("app's qpd exceeds limit, current: %d, limit: %d", appFlowCount.TotalCount, app.Qpd)).Error()))
		// 	c.Abort()
		// 	return
		// }

		c.Next()
	}
}
