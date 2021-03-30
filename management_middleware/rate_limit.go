package management_middleware

import (
	"github.com/LotteWong/giotto-gateway/common_middleware"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/time/rate"
	"net/http"
)

func RateLimitMiddleware() gin.HandlerFunc {
	rateFloat := lib.GetFloat64Conf("base.rate_limiter.rate")
	burstInt := lib.GetIntConf("base.rate_limiter.burst")
	limiter := rate.NewLimiter(rate.Limit(rateFloat), burstInt)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			common_middleware.ResponseError(c, http.StatusServiceUnavailable, errors.New("rate exceeds limit"))
			c.Abort()
			return
		}
		c.Next()
	}
}
