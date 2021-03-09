package po

type ServiceDetail struct {
	Info          *ServiceInfo   `json:"info" description:"基本信息"`
	TcpRule       *TcpRule       `json:"tcp_rule" desciption:"tcp_rule"`
	HttpRule      *HttpRule      `json:"http_rule" desciption:"http_rule"`
	GrpcRule      *GrpcRule      `json:"grpc_rule" desciption:"grpc_rule"`
	LoadBalance   *LoadBalance   `json:"load_balance" desciption:"load_balance"`
	AccessControl *AccessControl `json:"access_control" desciption:"access_control"`
}
