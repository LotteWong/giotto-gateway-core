package middleware

import (
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"log"
	"net"
	"net/http"
)

func CircuitBreakMiddleware() gin.HandlerFunc {
	// config hystrix circuit breaker
	timeout := lib.GetIntConf("base.circuit_breaker.timeout")
	maxConcurrentRequests := lib.GetIntConf("base.circuit_breaker.max_concurrent_request")
	sleepWindow := lib.GetIntConf("base.circuit_breaker.sleep_window")
	requestVolumeThreshold := lib.GetIntConf("base.circuit_breaker.request_volume_threshold")
	errorPercentThreshold := lib.GetIntConf("base.circuit_breaker.error_percent_threshold")
	hystrix.ConfigureCommand("gateway_hystrix", hystrix.CommandConfig{
		Timeout:                timeout,
		MaxConcurrentRequests:  maxConcurrentRequests,
		SleepWindow:            sleepWindow,
		RequestVolumeThreshold: requestVolumeThreshold,
		ErrorPercentThreshold:  errorPercentThreshold,
	})

	// config hystrix stream handler
	enableDashboard := lib.GetBoolConf("base.circuit_breaker.enable_dashboard")
	dashboardPort := lib.GetStringConf("base.circuit_breaker.dashboard_port")
	if enableDashboard {
		hystrixStreamHandler := hystrix.NewStreamHandler()
		hystrixStreamHandler.Start()
		go func() {
			err := http.ListenAndServe(net.JoinHostPort("", dashboardPort), hystrixStreamHandler)
			log.Fatal(err)
		}()
	}

	return func(c *gin.Context) {
		err := hystrix.Do("gateway_hystrix", func() error {
			c.Next()

			statusCode := c.Writer.Status()
			if !(statusCode >= 200 && statusCode < 300) {
				return errors.New("circuit breaker downstream error")
			}
			return nil
		}, nil)
		if err != nil {
			// TODO: service downgrade

			switch err {
			case hystrix.ErrCircuitOpen:
				ResponseError(c, http.StatusInternalServerError, errors.New(fmt.Sprintf("circuit open error:"+err.Error())))
			case hystrix.ErrMaxConcurrency:
				ResponseError(c, http.StatusInternalServerError, errors.New(fmt.Sprintf("circuit max concurrency error:"+err.Error())))
			default:
				ResponseError(c, http.StatusInternalServerError, errors.New(fmt.Sprintf("circuit default error:"+err.Error())))
			}
			c.Abort()
			return
		}
	}
}
