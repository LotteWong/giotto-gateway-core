package tcp_proxy_router

import (
	tcp_server "github.com/LotteWong/tcp-conn-server"
)

func InitRouter(middlewares ...tcp_server.TCPHandler) *tcp_server.TcpRouter {
	router := tcp_server.NewTcpRouter()
	router.Group("/").Use()
	return router
}
