package service

import (
	"github.com/LotteWong/giotto-gateway/dao/redis"
	"github.com/LotteWong/giotto-gateway/models/po"
	"sync"
	"time"
)

var flowCountService *FlowCountService

type FlowCountService struct {
	FlowCountMap   map[string]*po.FlowCount
	FlowCountSlice []*po.FlowCount
	RWLock         sync.RWMutex

	flowCountOperator *redis.FlowCountOperator
}

func NewFlowCountService() *FlowCountService {
	service := &FlowCountService{
		FlowCountMap:      map[string]*po.FlowCount{},
		FlowCountSlice:    []*po.FlowCount{},
		RWLock:            sync.RWMutex{},
		flowCountOperator: redis.NewFlowCountOperator(),
	}
	return service
}

func GetFlowCountService() *FlowCountService {
	if flowCountService == nil {
		flowCountService = NewFlowCountService()
	}
	return flowCountService
}

func (s *FlowCountService) GetFlowCount(serviceName string) (*po.FlowCount, error) {
	// hit in cache, use cache data
	for _, count := range s.FlowCountSlice {
		if count.ServiceName == serviceName {
			return count, nil
		}
	}

	// miss in cache, new a flow count
	req := &po.FlowCount{
		ServiceName: serviceName,
		Interval:    1 * time.Second,
	}
	count := s.flowCountOperator.GetFlowCount(req)

	// miss in cache, write back to cache
	s.FlowCountSlice = append(s.FlowCountSlice, count)
	s.RWLock.Lock()
	defer s.RWLock.Unlock()
	s.FlowCountMap[serviceName] = count

	return count, nil
}

func (s *FlowCountService) Increase(req *po.FlowCount) {
	s.flowCountOperator.Increase(req)
}

func (s *FlowCountService) GetDayFlow(t time.Time, serviceName string) (int64, error) {
	return s.flowCountOperator.GetRedisDayVal(t, &po.FlowCount{ServiceName: serviceName})
}

func (s *FlowCountService) GetHourFlow(t time.Time, serviceName string) (int64, error) {
	return s.flowCountOperator.GetRedisHourVal(t, &po.FlowCount{ServiceName: serviceName})
}
