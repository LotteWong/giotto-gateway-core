package service

import (
	"sync"

	"github.com/LotteWong/giotto-gateway-core/models/po"
	"github.com/go-redis/redis"
	"github.com/go-redis/redis_rate"
)

var rateLimitService *RateLimitService

type RateLimitService struct {
	RateLimitMap   map[string]*po.RateLimit
	RateLimitSlice []*po.RateLimit
	RWLock         sync.RWMutex
}

func NewRateLimitService() *RateLimitService {
	service := &RateLimitService{
		RateLimitMap:   map[string]*po.RateLimit{},
		RateLimitSlice: []*po.RateLimit{},
		RWLock:         sync.RWMutex{},
	}
	return service
}

func GetRateLimitService() *RateLimitService {
	if rateLimitService == nil {
		rateLimitService = NewRateLimitService()
	}
	return rateLimitService
}

// func (s *RateLimitService) GetRateLimit(serviceName string, qps int64) (*rate.Limiter, error) {
func (s *RateLimitService) GetRateLimit(serviceName string) (*redis_rate.Limiter, error) {
	// hit in cache, use cache data
	for _, limit := range s.RateLimitSlice {
		if limit.ServiceName == serviceName {
			return limit.Limiter, nil
		}
	}

	// miss in cache, new a rate limit
	// limiter := rate.NewLimiter(rate.Limit(qps), 3*int(qps))
	// TODO: use config from file system
	conn := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	limiter := redis_rate.NewLimiter(conn)

	// miss in cache, write back to cache
	limit := &po.RateLimit{
		ServiceName: serviceName,
		Limiter:     limiter,
	}
	s.RateLimitSlice = append(s.RateLimitSlice, limit)
	s.RWLock.Lock()
	defer s.RWLock.Unlock()
	s.RateLimitMap[serviceName] = limit

	return limiter, nil
}
