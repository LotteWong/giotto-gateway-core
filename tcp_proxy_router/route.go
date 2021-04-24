package tcp_proxy_router

import (
	"github.com/LotteWong/giotto-gateway-core/tcp_proxy_middleware"
	tcp "github.com/LotteWong/tcp-conn-server"
)

func InitRouter(middlewares ...tcp.TCPHandler) *tcp.TcpRouter {
	router := tcp.NewTcpRouter()
	router.Group("/").Use(
		tcp_proxy_middleware.TcpFlowCountMiddleware(),
		tcp_proxy_middleware.TcpRateLimitMiddleware(),
		tcp_proxy_middleware.TcpIpAuthMiddleware(),
	)
	return router
}
