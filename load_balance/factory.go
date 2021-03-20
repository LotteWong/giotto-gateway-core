package load_balance

import (
	"github.com/LotteWong/giotto-gateway/constants"
	"github.com/LotteWong/giotto-gateway/load_balance/lb_algos"
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
		return lb_algos.NewRandomLb()
	case constants.LbTypeRoundRobin:
		return lb_algos.NewRoundRobinLb()
	case constants.LbTypeWeightRoundRobin:
		return lb_algos.NewWeightRoundRobinLb()
	case constants.LbTypeConsistentHash:
		return lb_algos.NewConsistentHashLb(constants.DefaultReplicas, constants.DefaultHashFunc)
	default:
		return lb_algos.NewRandomLb()
	}
}

func LoadBalanceWithConfFactory(lbType LbType, conf LoadBalanceConf) LoadBalance {
	switch lbType {
	case constants.LbTypeRandom:
		lb := lb_algos.NewRandomLb()
		lb.Register(conf)
		lb.Subscribe()
		return lb
	case constants.LbTypeRoundRobin:
		lb := lb_algos.NewRoundRobinLb()
		lb.Register(conf)
		lb.Subscribe()
		return lb
	case constants.LbTypeWeightRoundRobin:
		lb := lb_algos.NewWeightRoundRobinLb()
		lb.Register(conf)
		lb.Subscribe()
		return lb
	case constants.LbTypeConsistentHash:
		lb := lb_algos.NewConsistentHashLb(constants.DefaultReplicas, constants.DefaultHashFunc)
		lb.Register(conf)
		lb.Subscribe()
		return lb
	default:
		lb := lb_algos.NewRandomLb()
		lb.Register(conf)
		lb.Subscribe()
		return lb
	}
}
