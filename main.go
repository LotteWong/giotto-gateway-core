package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	http_proxy_router "github.com/LotteWong/giotto-gateway-core/http_proxy_router"
	"github.com/LotteWong/giotto-gateway-core/service"
	"github.com/e421083458/golang_common/lib"
)

var (
	config = flag.String("config", "", "config file path")
)

func main() {
	flag.Parse()
	if *config == "" {
		flag.Usage()
		os.Exit(1)
	}

	InitCoreServer(*config)
}

func InitCoreServer(config string) {
	lib.InitModule(config, []string{"base", "mysql", "redis"})
	defer lib.Destroy()

	// _ = service.GetSvcService().LoadServicesIntoMemory()
	_ = service.GetSvcService().LoadServicesFromRedis()
	// _ = service.GetAppService().LoadAppsIntoMemory()
	_ = service.GetAppService().LoadAppsFromRedis()

	go func() {
		http_proxy_router.HttpServerRun()
	}()
	// go func() {
	// 	http_proxy_router.HttpsServerRun()
	// }()
	// go func() {
	// 	tcp_proxy_router.TcpServerRun()
	// }()
	// go func() {
	// 	grpc_proxy_router.GrpcServerRun()
	// }()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// grpc_proxy_router.GrpcServerStop()
	// tcp_proxy_router.TcpServerStop()
	// http_proxy_router.HttpsServerStop()
	http_proxy_router.HttpServerStop()
}
