package tcp_proxy_router

import (
	"context"
	"fmt"
	"log"

	"github.com/LotteWong/giotto-gateway/models/po"
	"github.com/LotteWong/giotto-gateway/reverse_proxy"
	"github.com/LotteWong/giotto-gateway/service"
	tcp "github.com/LotteWong/tcp-conn-server"
)

var tcpServers []*tcp.TcpServer

func TcpServerRun() {
	_, tcpServices, _, err := service.GetSvcService().GroupServicesInMemory()
	if err != nil {
		log.Fatalf(" [ERROR] TcpServerRun - err:%v\n", err)
	}

	for _, tcpService := range tcpServices {
		tmpService := tcpService
		go func(serviceDetail *po.ServiceDetail) {
			addr := fmt.Sprintf(":%d", serviceDetail.TcpRule.Port)

			lb, err := service.GetLbService().GetLbWithConfForSvc(serviceDetail)
			if err != nil {
				log.Fatalf(" [ERROR] TcpServerRun - tcp proxy server:%s err:%v\n", addr, err)
			}
			r := InitRouter()

			tcpSrvHandler := tcp.NewTcpRouteHandler(r, func(ctx *tcp.TcpRouterContext) tcp.TCPHandler {
				return reverse_proxy.NewTcpReverseProxy(ctx, lb)
			})
			baseCtx := context.WithValue(context.Background(), "service", serviceDetail)

			tcpServer := &tcp.TcpServer{
				Addr:    addr,
				Handler: tcpSrvHandler,
				BaseCtx: baseCtx,
			}
			tcpServers = append(tcpServers, tcpServer)

			log.Printf(" [INFO] TcpServerRun - tcp proxy server:%s\n", addr)
			if err := tcpServer.ListenAndServe(); err != nil && err != tcp.ErrServerClosed {
				log.Fatalf(" [ERROR] TcpServerRun - tcp proxy server:%s err:%v\n", addr, err)
			}
		}(tmpService)
	}
}

func TcpServerStop() {
	for _, tcpServer := range tcpServers {
		if err := tcpServer.Close(); err != nil {
			log.Fatalf(" [ERROR] TcpServerStop - tcp proxy server err:%v\n", err)
		}
		log.Printf(" [INFO] TcpServerStop - tcp proxy server stopped\n")
	}
}
