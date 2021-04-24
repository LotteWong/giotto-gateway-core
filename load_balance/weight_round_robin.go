package load_balance

import (
	"strconv"

	"github.com/pkg/errors"
)

type WeightRoundRobinLb struct {
	nodes []*IpAndWeightNode
	idx   int
	conf  LoadBalanceConf
}

type IpAndWeightNode struct {
	ip              string
	weight          int
	currentWeight   int // may greater than weight
	effectiveWeight int // must less than or equal to weight
}

func NewWeightRoundRobinLb() *WeightRoundRobinLb {
	return &WeightRoundRobinLb{
		nodes: []*IpAndWeightNode{},
		idx:   0,
		conf:  nil,
	}
}

func (lb *WeightRoundRobinLb) Add(params ...string) error {
	if len(params) != 2 {
		return errors.New("params length should be 2")
	}
	ip := params[0]
	weight, err := strconv.Atoi(params[1])
	if err != nil {
		return err
	}

	lb.nodes = append(lb.nodes, &IpAndWeightNode{
		ip:              ip,
		weight:          weight,
		effectiveWeight: weight,
	})

	return nil
}

func (lb *WeightRoundRobinLb) Rmv(params ...string) error {
	if len(params) == 0 {
		return errors.New("params length is at least 1")
	}
	ip := params[0]

	var newNodes []*IpAndWeightNode
	for _, oldNode := range lb.nodes {
		if oldNode.ip == ip {
			continue
		}
		newNodes = append(newNodes, oldNode)
	}
	lb.nodes = newNodes

	return nil
}

func (lb *WeightRoundRobinLb) Get(key string) (string, error) {
	var bestNode *IpAndWeightNode
	var totalWeight int

	for _, node := range lb.nodes {
		totalWeight += node.effectiveWeight

		node.currentWeight += node.effectiveWeight
		if node.effectiveWeight < node.weight {
			node.effectiveWeight++
		}

		if bestNode == nil || node.currentWeight > bestNode.currentWeight {
			bestNode = node
		}
	}

	if bestNode == nil {
		return "", errors.New("no available ip")
	} else {
		bestNode.currentWeight -= totalWeight
		return bestNode.ip, nil
	}
}

func (lb *WeightRoundRobinLb) Register(conf LoadBalanceConf) {
	lb.conf = conf
	lb.conf.Attach(lb)
}

func (lb *WeightRoundRobinLb) Subscribe() {
	// TODO: strategy pattern improvement
	// if conf, ok := lb.conf.(*ClientSvcDiscoveryLbConf); ok {
	// 	lb.nodes = []*IpAndWeightNode{}
	// 	for _, pair := range conf.GetConf() {
	// 		lb.Add(pair.Ip, strconv.Itoa(pair.Weight))
	// 	}
	// }
	// if conf, ok := lb.conf.(*ServerSvcDiscoveryLbConf); ok {
	// 	lb.nodes = []*IpAndWeightNode{}
	// 	for _, pair := range conf.GetConf() {
	// 		lb.Add(pair.Ip, strconv.Itoa(pair.Weight))
	// 	}
	// }

	lb.nodes = []*IpAndWeightNode{}
	for _, pair := range lb.conf.GetConf() {
		lb.Add(pair.Ip, strconv.Itoa(pair.Weight))
	}
}
