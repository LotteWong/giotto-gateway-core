package po

import (
	"github.com/go-redis/redis_rate"
)

type RateLimit struct {
	ServiceName string `description:"服务名"`
	// Limiter     *rate.Limiter `description:"限流器"`
	Limiter *redis_rate.Limiter `description:"限流器"`
}
