package lb_conf

import (
	"github.com/LotteWong/giotto-gateway/constants"
	"github.com/LotteWong/giotto-gateway/load_balance"
	"net"
	"reflect"
	"sort"
	"time"
)

type ClientSvcDiscoveryLbConf struct {
	// observers
	lbs []load_balance.LoadBalance
	// configs
	activeIps   []string
	ipWeightMap map[string]int
	// others
	format string
}

func NewClientSvcDiscoveryLbConf(activeIps []string, ipWeightMap map[string]int, format string) *ClientSvcDiscoveryLbConf {
	// initiate conf
	conf := &ClientSvcDiscoveryLbConf{
		lbs:         []load_balance.LoadBalance{},
		ipWeightMap: ipWeightMap,
		activeIps:   activeIps,
		format:      format,
	}
	// publish conf
	conf.Publish()

	return conf
}

// Attach is for subject to attach observer
func (c *ClientSvcDiscoveryLbConf) Attach(lb load_balance.LoadBalance) {
	c.lbs = append(c.lbs, lb)
}

func (c *ClientSvcDiscoveryLbConf) Notify() {
	for _, lb := range c.lbs {
		lb.Subscribe()
	}
}

// Publish is for subject to publish to observer
func (c *ClientSvcDiscoveryLbConf) Publish() {
	// TODO: use job instead of loop
	go func() {
		connErrMap := map[string]int{}
		for {
			var newActiveIps []string

			// gateway client health check backend server
			for ip, _ := range c.ipWeightMap {
				conn, err := net.DialTimeout("tcp", ip, time.Duration(constants.DefaultDialTimeout)*time.Second)

				if err == nil {
					// if dial succeed, set ip's connection error nums to 0
					conn.Close()
					connErrMap[ip] = 0
				} else {
					// if dial failed, increase ip's connection error nums
					if _, ok := connErrMap[ip]; ok {
						connErrMap[ip] += 1
					} else {
						connErrMap[ip] = 1
					}
				}

				if connErrMap[ip] < constants.DefaultDialMaxErrNum {
					newActiveIps = append(newActiveIps, ip)
				}
			}

			// update active ips and Notify load balances
			sort.Strings(newActiveIps)
			sort.Strings(c.activeIps)
			if !reflect.DeepEqual(newActiveIps, c.activeIps) {
				c.activeIps = newActiveIps
				c.Notify()
			}

			time.Sleep(time.Duration(constants.DefaultDialInterval) * time.Second)
		}
	}()
}

func (c *ClientSvcDiscoveryLbConf) GetConf() []*load_balance.IpAndWeight {
	var confs []*load_balance.IpAndWeight

	for _, ip := range c.activeIps {
		weight, ok := c.ipWeightMap[ip]
		if !ok {
			weight = 50 // set default weight to 50
		}
		confs = append(confs, &load_balance.IpAndWeight{
			Ip:     ip,
			Weight: weight,
		})
	}

	return confs
}
