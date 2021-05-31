package constants

const (
	ServiceTypeHttp = 0
	ServiceTypeTcp  = 1
	ServiceTypeGrpc = 2

	HttpServiceTypeHttp  = 0
	HttpServiceTypeHttps = 1
	HttpServiceTypeWs    = 2
	HttpServiceTypeWss   = 3

	HttpRuleTypePrefixUrl = 0
	HttpRuleTypeDomain    = 1

	Disable = 0
	Enable  = 1

	LbTypeRandom           = 0
	LbTypeRoundRobin       = 1
	LbTypeWeightRoundRobin = 2
	LbTypeConsistentHash   = 3
)

var (
	ServiceTypeMap = map[int]string{
		ServiceTypeHttp: "Http",
		ServiceTypeTcp:  "Tcp",
		ServiceTypeGrpc: "Grpc",
	}

	HttpServiceTypeMap = map[int]string{
		HttpServiceTypeHttp:  "Http",
		HttpServiceTypeHttps: "Https",
		HttpServiceTypeWs:    "Ws",
		HttpServiceTypeWss:   "Wss",
	}
)
