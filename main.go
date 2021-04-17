package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	http_proxy_router "github.com/LotteWong/giotto-gateway/http_proxy_router"
	"github.com/LotteWong/giotto-gateway/management_router"
	"github.com/LotteWong/giotto-gateway/service"
	tcp_proxy_router "github.com/LotteWong/giotto-gateway/tcp_proxy_router"
	"github.com/e421083458/golang_common/lib"
)

var (
	endpoint = flag.String("endpoint", "", "management or proxy")
	config   = flag.String("config", "", "config file path")
)

func main() {
	flag.Parse()
	if *endpoint == "" || *config == "" {
		flag.Usage()
		os.Exit(1)
	}

	switch *endpoint {
	case "management":
		InitManagementServer(*config)
	case "proxy":
		InitProxyServer(*config)
	default:
		log.Fatalf(" [ERROR] endpoint %s is invalid", *endpoint)
	}
}

func InitManagementServer(config string) {
	lib.InitModule(config, []string{"base", "mysql", "redis"})
	defer lib.Destroy()

	management_router.HttpServerRun()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	management_router.HttpServerStop()
}

func InitProxyServer(config string) {
	lib.InitModule(config, []string{"base", "mysql", "redis"})
	defer lib.Destroy()
	_ = service.GetSvcService().LoadServicesIntoMemory()
	_ = service.GetAppService().LoadAppsIntoMemory()

	go func() {
		http_proxy_router.HttpServerRun()
	}()
	// go func() {
	// 	http_proxy_router.HttpsServerRun()
	// }()
	go func() {
		tcp_proxy_router.TcpServerRun()
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	tcp_proxy_router.TcpServerStop()
	http_proxy_router.HttpServerStop()
	// http_proxy_router.HttpsServerStop()
}
