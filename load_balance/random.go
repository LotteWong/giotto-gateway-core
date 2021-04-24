package load_balance

import (
	"math/rand"
	"strconv"

	"github.com/pkg/errors"
)

type RandomLb struct {
	ips  []string
	idx  int
	conf LoadBalanceConf
}

func NewRandomLb() *RandomLb {
	return &RandomLb{
		ips:  []string{},
		idx:  0,
		conf: nil,
	}
}

func (lb *RandomLb) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("params length is at least 1")
	}
	ip := params[0]

	lb.ips = append(lb.ips, ip)
	return nil
}

func (lb *RandomLb) Rmv(params ...string) error {
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

func (lb *RandomLb) Get(key string) (string, error) {
	if len(lb.ips) == 0 {
		return "", errors.New("no available ip")
	}

	lb.idx = rand.Intn(len(lb.ips))
	ip := lb.ips[lb.idx]

	return ip, nil
}

func (lb *RandomLb) Register(conf LoadBalanceConf) {
	lb.conf = conf
	lb.conf.Attach(lb)
}

// Subscribe is for observer to subscribe from subject
func (lb *RandomLb) Subscribe() {
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
