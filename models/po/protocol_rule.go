package po

type HttpRule struct {
	Id              int64  `json:"id" gorm:"primary_key" description:"自增主键"`
	ServiceId       int64  `json:"service_id" gorm:"column:service_id" description:"服务id"`
	RuleType        int    `json:"rule_type" gorm:"column:rule_type" description:"匹配类型 domain=域名，url_prefix=url前缀"`
	Rule            string `json:"rule" gorm:"column:rule" description:"type=domain表示域名，type=url_prefix表示url前缀"`
	NeedHttps       int    `json:"need_https" gorm:"column:need_https" description:"支持https 1=支持"`
	NeedStripUri    int    `json:"need_strip_uri" gorm:"column:need_strip_uri" description:"启用strip_uri 1=启用"`
	NeedWebsocket   int    `json:"need_websocket" gorm:"column:need_websocket" description:"支持websocket 1=支持"`
	UrlRewrite      string `json:"url_rewrite" gorm:"column:url_rewrite" description:"url重写功能 格式：^/gatekeeper/test_service(.*) $1 多个逗号间隔"`
	HeaderTransform string `json:"header_transform" gorm:"column:header_transform" description:"header转换支持增加(add)、删除(del)、修改(edit) 格式: add header_name header_value"`
}

func (t *HttpRule) TableName() string {
	return "gateway_service_http_rule"
}

type TcpRule struct {
	Id        int64 `json:"id" gorm:"primary_key" description:"自增主键"`
	ServiceId int64 `json:"service_id" gorm:"column:service_id" description:"服务id"`
	Port      int   `json:"port" gorm:"column:port" description:"端口"`
}

func (t *TcpRule) TableName() string {
	return "gateway_service_tcp_rule"
}

type GrpcRule struct {
	Id              int64  `json:"id" gorm:"primary_key" description:"自增主键"`
	ServiceId       int64  `json:"service_id" gorm:"column:service_id" description:"服务id	"`
	Port            int    `json:"port" gorm:"column:port" description:"端口"`
	HeaderTransform string `json:"header_transform" gorm:"column:header_transform" description:"header转换支持增加(add)、删除(del)、修改(edit) 格式: add header_name header_value"`
}

func (t *GrpcRule) TableName() string {
	return "gateway_service_grpc_rule"
}
