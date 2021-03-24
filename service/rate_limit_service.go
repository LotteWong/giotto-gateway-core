package service

import (
	"github.com/LotteWong/giotto-gateway/models/po"
	"golang.org/x/time/rate"
	"sync"
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

func (s *RateLimitService) GetRateLimit(serviceName string, qps int64) (*rate.Limiter, error) {
	// hit in cache, use cache data
	for _, limit := range s.RateLimitSlice {
		if limit.ServiceName == serviceName {
			return limit.Limiter, nil
		}
	}

	// miss in cache, new a rate limit
	limiter := rate.NewLimiter(rate.Limit(qps), 3*int(qps))

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
