package reverse_proxy

import (
	"github.com/LotteWong/giotto-gateway/load_balance"
	"github.com/LotteWong/giotto-gateway/middleware"
	"github.com/LotteWong/giotto-gateway/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func NewReverseProxy(ctx *gin.Context, lb load_balance.LoadBalance, trans *http.Transport) *httputil.ReverseProxy {
	// convert the source request to target request
	director := func(req *http.Request) {
		var err error

		targetStr, err := lb.Get(req.URL.String())
		if err != nil {
			panic(err)
		}
		targetUrl, err := url.Parse(targetStr)
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
		if strings.Contains(res.Header.Get("Connection"), "Upgrade") {
			// TODO
		}
		return nil
	}

	// handle when error occurs
	errFunc := func(w http.ResponseWriter, r *http.Request, err error) {
		middleware.ResponseError(ctx, http.StatusInternalServerError, err)
	}

	return &httputil.ReverseProxy{
		Transport:      trans,
		Director:       director,
		ModifyResponse: modifyFunc,
		ErrorHandler:   errFunc,
	}
}