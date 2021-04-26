package router

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/LotteWong/giotto-gateway-core/common_middleware"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
)

var (
	HttpSrvHandler  *http.Server
	HttpsSrvHandler *http.Server
)

func HttpServerRun() {
	gin.SetMode(lib.ConfBase.DebugMode)
	r := InitRouter([]gin.HandlerFunc{
		common_middleware.RecoveryMiddleware(),
		common_middleware.RequestLog(),
	}...)
	HttpSrvHandler = &http.Server{
		Addr:           lib.GetStringConf("base.http.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(lib.GetIntConf("base.http.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(lib.GetIntConf("base.http.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(lib.GetIntConf("base.http.max_header_bytes")),
	}
	log.Printf(" [INFO] HttpServerRun - http proxy server:%s\n", lib.GetStringConf("base.http.addr"))
	if err := HttpSrvHandler.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf(" [ERROR] HttpServerRun - http proxy server:%s err:%v\n", lib.GetStringConf("base.http.addr"), err)
	}
}

func HttpServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpSrvHandler.Shutdown(ctx); err != nil {
		log.Fatalf(" [ERROR] HttpServerStop - http proxy server err:%v\n", err)
	}
	log.Printf(" [INFO] HttpServerStop - http proxy server stopped\n")
}

func HttpsServerRun() {
	gin.SetMode(lib.ConfBase.DebugMode)
	r := InitRouter([]gin.HandlerFunc{
		common_middleware.RecoveryMiddleware(),
		common_middleware.RequestLog(),
	}...)
	HttpsSrvHandler = &http.Server{
		Addr:           lib.GetStringConf("base.https.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(lib.GetIntConf("base.https.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(lib.GetIntConf("base.https.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(lib.GetIntConf("base.https.max_header_bytes")),
	}
	log.Printf(" [INFO] HttpsServerRun - https proxy server:%s\n", lib.GetStringConf("base.https.addr"))
	certFile := lib.GetStringConf("base.https.cert_file")
	keyFile := lib.GetStringConf("base.https.key_file")
	if err := HttpsSrvHandler.ListenAndServeTLS(certFile, keyFile); err != nil && err != http.ErrServerClosed {
		log.Fatalf(" [ERROR] HttpsServerRun - https proxy server:%s err:%v\n", lib.GetStringConf("base.https.addr"), err)
	}
}

func HttpsServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpsSrvHandler.Shutdown(ctx); err != nil {
		log.Fatalf(" [ERROR] HttpsServerStop - https proxy server err:%v\n", err)
	}
	log.Printf(" [INFO] HttpsServerStop - https proxy server stopped\n")
}
