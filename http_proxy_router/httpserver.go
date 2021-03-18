package router

import (
	"context"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

var (
	HttpSrvHandler  *http.Server
	HttpsSrvHandler *http.Server
)

func HttpServerRun() {
	gin.SetMode(lib.GetStringConf("proxy.base.debug_mode"))
	r := InitRouter()
	HttpSrvHandler = &http.Server{
		Addr:           lib.GetStringConf("proxy.http.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(lib.GetIntConf("proxy.http.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(lib.GetIntConf("proxy.http.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(lib.GetIntConf("proxy.http.max_header_bytes")),
	}
	log.Printf(" [INFO] HttpServerRun - http proxy server:%s\n", lib.GetStringConf("proxy.http.addr"))
	if err := HttpSrvHandler.ListenAndServe(); err != nil {
		log.Fatalf(" [ERROR] HttpServerRun - http proxy server:%s err:%v\n", lib.GetStringConf("proxy.http.addr"), err)
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
	gin.SetMode(lib.GetStringConf("proxy.base.debug_mode"))
	r := InitRouter()
	HttpsSrvHandler = &http.Server{
		Addr:           lib.GetStringConf("proxy.https.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(lib.GetIntConf("proxy.https.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(lib.GetIntConf("proxy.https.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(lib.GetIntConf("proxy.https.max_header_bytes")),
	}
	log.Printf(" [INFO] HttpsServerRun - https proxy server:%s\n", lib.GetStringConf("proxy.https.addr"))
	certFile := lib.GetStringConf("proxy.https.cert_file")
	keyFile := lib.GetStringConf("proxy.https.key_file")
	if err := HttpsSrvHandler.ListenAndServeTLS(certFile, keyFile); err != nil {
		log.Fatalf(" [ERROR] HttpsServerRun - https proxy server:%s err:%v\n", lib.GetStringConf("proxy.https.addr"), err)
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
