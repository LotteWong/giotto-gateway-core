package po

import "golang.org/x/time/rate"

type RateLimit struct {
	ServiceName string        `description:"服务名"`
	Limiter     *rate.Limiter `description:"限流器"`
}
