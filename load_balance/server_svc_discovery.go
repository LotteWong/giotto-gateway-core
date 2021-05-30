package load_balance

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"

	"github.com/e421083458/golang_common/lib"
	"github.com/hashicorp/consul/api"
)

type ServerSvcDiscoveryLbConf struct {
	// observers
	lbs []LoadBalance
	// configs
	activeIps   []string
	ipWeightMap map[string]int
	// others
	format  string
	service string
	tag     string
}

func NewServerSvcDiscoveryLbConf(activeIps []string, ipWeightMap map[string]int, format, service, tag string) *ServerSvcDiscoveryLbConf {
	// initiate conf
	conf := &ServerSvcDiscoveryLbConf{
		lbs:         []LoadBalance{},
		ipWeightMap: ipWeightMap,
		activeIps:   activeIps,
		format:      format,
		service:     service,
		tag:         tag,
	}
	// publish conf
	conf.Publish()

	return conf
}

// Attach is for subject to attach observer
func (c *ServerSvcDiscoveryLbConf) Attach(lb LoadBalance) {
	c.lbs = append(c.lbs, lb)
}

func (c *ServerSvcDiscoveryLbConf) Notify() {
	for _, lb := range c.lbs {
		lb.Subscribe()
	}
}

// Publish is for subject to publish to observer
func (c *ServerSvcDiscoveryLbConf) Publish() {
	consulIp := lib.GetStringConf("base.consul.ip")
	consulPort := lib.GetIntConf("base.consul.port")

	// init consul client
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%d", consulIp, consulPort)
	client, err := api.NewClient(config)
	if err != nil {
		log.Printf("failed to new a consul client, err: %v\n", err)
		return
	}

	// TODO: use job instead of loop
	go func() {
		var lastIndex uint64
		var isConfigured bool = len(c.ipWeightMap) > 0

		for {
			var newActiveIps []string

			// gateway server health check backend server
			serviceEntries, meta, err := client.Health().Service(c.service, "", true, &api.QueryOptions{
				WaitIndex: lastIndex,
			})
			if err != nil {
				log.Printf("failed to get services from consul, err: %v\n", err)
			}
			lastIndex = meta.LastIndex
			for _, serviceEntry := range serviceEntries {
				ip := fmt.Sprintf("%s:%d", serviceEntry.Service.Address, serviceEntry.Service.Port)
				weight, _ := strconv.Atoi(serviceEntry.Service.Meta["weight"])

				// if ip list and weight list are configured, follow the load balance configs
				if isConfigured {
					if _, ok := c.ipWeightMap[ip]; ok {
						newActiveIps = append(newActiveIps, ip)
					}
				}

				// if ip list and weight list are not configured, follow the service discovery configs
				if !isConfigured {
					c.ipWeightMap[ip] = weight
					newActiveIps = append(newActiveIps, ip)
				}
			}
			log.Printf("consul watch ips update: %v", newActiveIps)

			// update active ips and Notify load balances
			sort.Strings(newActiveIps)
			sort.Strings(c.activeIps)
			if !reflect.DeepEqual(newActiveIps, c.activeIps) {
				c.activeIps = newActiveIps
				c.Notify()
			}
		}
	}()
}

func (c *ServerSvcDiscoveryLbConf) GetConf() []*IpAndWeight {
	var confs []*IpAndWeight

	for _, ip := range c.activeIps {
		weight, ok := c.ipWeightMap[ip]
		if !ok {
			weight = 50 // set default weight to 50
		}
		confs = append(confs, &IpAndWeight{
			Ip:     ip,
			Weight: weight,
		})
	}

	return confs
}
