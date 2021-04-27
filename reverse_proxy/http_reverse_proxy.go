package reverse_proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/LotteWong/giotto-gateway-core/common_middleware"
	"github.com/LotteWong/giotto-gateway-core/load_balance"
	"github.com/LotteWong/giotto-gateway-core/utils"
	"github.com/gin-gonic/gin"
)

func NewHttpReverseProxy(ctx *gin.Context, lb load_balance.LoadBalance, trans *http.Transport, scheme string) *httputil.ReverseProxy {
	// convert the source request to target request
	director := func(req *http.Request) {
		var err error

		targetStr, err := lb.Get(req.URL.String())
		if err != nil {
			panic(err)
		}
		targetUrl, err := url.Parse(scheme + targetStr)
		if err != nil {
			panic(err)
		}

		req.URL.Scheme = targetUrl.Scheme
		req.URL.Host = targetUrl.Host
		req.URL.Path = utils.JoinSlash(targetUrl.Path, req.URL.Path)
		req.Host = targetUrl.Host
		if targetUrl.RawQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = req.URL.RawQuery + targetUrl.RawQuery
		} else {
			req.URL.RawQuery = req.URL.RawQuery + "&" + targetUrl.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "user-agent")
		}
	}

	// modify the source response to target response
	modifyFunc := func(res *http.Response) error {
		// support websocket protocol
		if strings.Contains(res.Header.Get("Connection"), "Upgrade") {
			return nil
		}
		return nil
	}

	// handle when error occurs
	errFunc := func(w http.ResponseWriter, r *http.Request, err error) {
		common_middleware.ResponseError(ctx, http.StatusInternalServerError, err)
	}

	return &httputil.ReverseProxy{
		Transport:      trans,
		Director:       director,
		ModifyResponse: modifyFunc,
		ErrorHandler:   errFunc,
	}
}
