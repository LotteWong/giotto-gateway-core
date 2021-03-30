package load_balance

import (
	"github.com/LotteWong/giotto-gateway/constants"
)

type LbType int

type LoadBalance interface {
	Add(...string) error
	Rmv(...string) error
	Get(string) (string, error)

	Subscribe()
}

type LoadBalanceConf interface {
	GetConf() []*IpAndWeight

	Attach(lb LoadBalance)
	Publish()
}

type IpAndWeight struct {
	Ip     string
	Weight int
}

func LoadBalanceFactory(lbType LbType) LoadBalance {
	switch lbType {
	case constants.LbTypeRandom:
		return NewRandomLb()
	case constants.LbTypeRoundRobin:
		return NewRoundRobinLb()
	case constants.LbTypeWeightRoundRobin:
		return NewWeightRoundRobinLb()
	case constants.LbTypeConsistentHash:
		return NewConsistentHashLb(DefaultReplicas, DefaultHashFunc)
	default:
		return NewRandomLb()
	}
}

func LoadBalanceWithConfFactory(lbType LbType, conf LoadBalanceConf) LoadBalance {
	switch lbType {
	case constants.LbTypeRandom:
		lb := NewRandomLb()
		lb.Register(conf)
		lb.Subscribe()
		return lb
	case constants.LbTypeRoundRobin:
		lb := NewRoundRobinLb()
		lb.Register(conf)
		lb.Subscribe()
		return lb
	case constants.LbTypeWeightRoundRobin:
		lb := NewWeightRoundRobinLb()
		lb.Register(conf)
		lb.Subscribe()
		return lb
	case constants.LbTypeConsistentHash:
		lb := NewConsistentHashLb(DefaultReplicas, DefaultHashFunc)
		lb.Register(conf)
		lb.Subscribe()
		return lb
	default:
		lb := NewRandomLb()
		lb.Register(conf)
		lb.Subscribe()
		return lb
	}
}
