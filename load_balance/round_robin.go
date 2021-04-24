package load_balance

import (
	"strconv"

	"github.com/pkg/errors"
)

type RoundRobinLb struct {
	ips  []string
	idx  int
	conf LoadBalanceConf
}

func NewRoundRobinLb() *RoundRobinLb {
	return &RoundRobinLb{
		ips:  []string{},
		idx:  0,
		conf: nil,
	}
}

func (lb *RoundRobinLb) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("params length is at least 1")
	}
	ip := params[0]

	lb.ips = append(lb.ips, ip)
	return nil
}

func (lb *RoundRobinLb) Rmv(params ...string) error {
	if len(params) == 0 {
		return errors.New("params length is at least 1")
	}
	ip := params[0]

	var newIps []string
	for _, oldIp := range lb.ips {
		if oldIp == ip {
			continue
		}
		newIps = append(newIps, oldIp)
	}
	lb.ips = newIps

	return nil
}

func (lb *RoundRobinLb) Get(key string) (string, error) {
	if len(lb.ips) == 0 {
		return "", errors.New("no available ip")
	}

	ipLen := len(lb.ips)
	if lb.idx >= ipLen {
		lb.idx = 0
	}
	ip := lb.ips[lb.idx]
	lb.idx = (lb.idx + 1) % ipLen

	return ip, nil
}

func (lb *RoundRobinLb) Register(conf LoadBalanceConf) {
	lb.conf = conf
	lb.conf.Attach(lb)
}

// Subscribe is for observer to subscribe from subject
func (lb *RoundRobinLb) Subscribe() {
	// TODO: strategy pattern improvement
	// if conf, ok := lb.conf.(*ClientSvcDiscoveryLbConf); ok {
	// 	lb.ips = []string{}
	// 	for _, pair := range conf.GetConf() {
	// 		lb.Add(pair.Ip, strconv.Itoa(pair.Weight))
	// 	}
	// }
	// if conf, ok := lb.conf.(*ServerSvcDiscoveryLbConf); ok {
	// 	lb.ips = []string{}
	// 	for _, pair := range conf.GetConf() {
	// 		lb.Add(pair.Ip, strconv.Itoa(pair.Weight))
	// 	}
	// }

	lb.ips = []string{}
	for _, pair := range lb.conf.GetConf() {
		lb.Add(pair.Ip, strconv.Itoa(pair.Weight))
	}
}
