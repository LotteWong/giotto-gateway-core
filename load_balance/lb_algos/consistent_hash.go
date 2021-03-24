package lb_algos

import (
	"github.com/LotteWong/giotto-gateway/load_balance"
	"github.com/LotteWong/giotto-gateway/load_balance/lb_conf"
	"github.com/pkg/errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type HashFunc func(data []byte) uint32

type HashRing []uint32

func (r HashRing) Len() int {
	return len(r)
}

func (r HashRing) Less(i, j int) bool {
	return r[i] < r[j]
}

func (r HashRing) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

type ConsistentHashLb struct {
	hashMap  map[uint32]string // key is hash, value is ip
	hashRing HashRing
	hashFunc HashFunc
	replicas int
	RWLock   sync.RWMutex
	conf     load_balance.LoadBalanceConf
}

func NewConsistentHashLb(replicas int, function HashFunc) *ConsistentHashLb {
	if replicas == 0 {
		replicas = 1 // set default replicas to 1
	}
	if function == nil {
		function = crc32.ChecksumIEEE // set default hash func as crc32
	}
	lb := &ConsistentHashLb{
		hashMap:  make(map[uint32]string),
		hashRing: HashRing{},
		hashFunc: function,
		replicas: replicas,
		RWLock:   sync.RWMutex{},
		conf:     nil,
	}
	return lb
}

func (lb *ConsistentHashLb) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("params length is at least 1")
	}
	ip := params[0]

	lb.RWLock.Lock()
	defer lb.RWLock.Unlock()

	for i := 0; i < lb.replicas; i++ {
		hash := lb.hashFunc([]byte(strconv.Itoa(i) + ip))
		lb.hashRing = append(lb.hashRing, hash)
		lb.hashMap[hash] = ip
	}
	sort.Sort(lb.hashRing) // easy to binary search

	return nil
}

func (lb *ConsistentHashLb) Rmv(params ...string) error {
	// TODO: remove hash node from hash ring
	return nil
}

func (lb *ConsistentHashLb) Get(key string) (string, error) {
	if len(lb.hashRing) == 0 {
		return "", errors.New("no available ip")
	}
	keyHash := lb.hashFunc([]byte(key))

	idx := sort.Search(len(lb.hashRing), func(i int) bool {
		return lb.hashRing[i] >= keyHash
	})
	if idx == len(lb.hashRing) {
		idx = 0
	}

	lb.RWLock.Lock()
	lb.RWLock.Unlock()

	hash := lb.hashRing[idx]
	ip := lb.hashMap[hash]
	return ip, nil
}

func (lb *ConsistentHashLb) Register(conf load_balance.LoadBalanceConf) {
	lb.conf = conf
	lb.conf.Attach(lb)
}

func (lb *ConsistentHashLb) Subscribe() {
	if conf, ok := lb.conf.(*lb_conf.ClientSvcDiscoveryLbConf); ok {
		lb.hashMap = map[uint32]string{}
		lb.hashRing = HashRing{}
		for _, pair := range conf.GetConf() {
			lb.Add(pair.Ip, strconv.Itoa(pair.Weight))
		}
	}
}
