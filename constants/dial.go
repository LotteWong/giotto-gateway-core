package constants

import "github.com/LotteWong/giotto-gateway/load_balance/lb_algos"

const (
	DefaultDialMethod    = 0
	DefaultDialTimeout   = 5
	DefaultDialMaxErrNum = 3
	DefaultDialInterval  = 5

	DefaultReplicas = 10
)

var (
	DefaultHashFunc lb_algos.HashFunc = nil
)
