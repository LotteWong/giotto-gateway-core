package po

import (
	"strings"

	"github.com/LotteWong/giotto-gateway-core/load_balance"
)

type LoadBalance struct {
	Id        int64 `json:"id" gorm:"primary_key" description:"自增主键"`
	ServiceId int64 `json:"service_id" gorm:"column:service_id" description:"服务id	"`

	CheckMethod   int    `json:"check_method" gorm:"column:check_method" description:"检查方法 tcpchk=检测端口是否握手成功"`
	CheckTimeout  int    `json:"check_timeout" gorm:"column:check_timeout" description:"超时时间，单位为s"`
	CheckInterval int    `json:"check_interval" gorm:"column:check_interval" description:"检查间隔，单位为s"`
	RoundType     int    `json:"round_type" gorm:"column:round_type" description:"轮询方式 round/weight_round/random/ip_hash"`
	IpList        string `json:"ip_list" gorm:"column:ip_list" description:"启用ip列表"`
	WeightList    string `json:"weight_list" gorm:"column:weight_list" description:"权重列表"`
	ForbidList    string `json:"forbid_list" gorm:"column:forbid_list" description:"禁用ip列表"`

	UpstreamConnectTimeout int `json:"upstream_connect_timeout" gorm:"column:upstream_connect_timeout" description:"下游建立连接超时, 单位为s"`
	UpstreamHeaderTimeout  int `json:"upstream_header_timeout" gorm:"column:upstream_header_timeout" description:"下游获取头部超时, 单位为s"`
	UpstreamIdleTimeout    int `json:"upstream_idle_timeout" gorm:"column:upstream_idle_timeout" description:"下游连接最大空闲时间, 单位为s"`
	UpstreamMaxIdle        int `json:"upstream_max_idle" gorm:"column:upstream_max_idle" description:"下游最大空闲连接数量"`

	IsDelete int `json:"is_delete" gorm:"column:is_delete" description:"是否删除"`
}

func (t *LoadBalance) TableName() string {
	return "gateway_service_load_balance"
}

func (t *LoadBalance) GetEnabledIpList() []string {
	return strings.Split(t.IpList, ",")
}

func (t *LoadBalance) GetDisabledIpList() []string {
	return strings.Split(t.ForbidList, ",")
}

func (t *LoadBalance) GetWeightList() []string {
	return strings.Split(t.WeightList, ",")
}

type LoadBalanceDetail struct {
	LoadBalancer load_balance.LoadBalance `json:"load_balance" description:"负载均衡器"`
	ServiceName  string                   `json:"service_name" description:"服务名称"`
}
