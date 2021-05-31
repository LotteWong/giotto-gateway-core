package po

import "time"

type ServiceInfo struct {
	Id          int64     `json:"id" gorm:"primary_key" description:"自增主键"`
	ServiceName string    `json:"service_name" gorm:"column:service_name" description:"服务名称"`
	ServiceDesc string    `json:"service_desc" gorm:"column:service_desc" description:"服务描述"`
	ServiceType int       `json:"service_type" gorm:"column:service_type" description:"服务类型"`
	UpdatedAt   time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	CreatedAt   time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	IsDelete    int       `json:"is_delete" gorm:"column:is_delete" description:"是否删除"`
}

func (t *ServiceInfo) TableName() string {
	return "gateway_service_info"
}

type ServiceDetail struct {
	Info          *ServiceInfo   `json:"info" description:"基本信息"`
	TcpRule       *TcpRule       `json:"tcp_rule" desciption:"tcp_rule"`
	HttpRule      *HttpRule      `json:"http_rule" desciption:"http_rule"`
	GrpcRule      *GrpcRule      `json:"grpc_rule" desciption:"grpc_rule"`
	LoadBalance   *LoadBalance   `json:"load_balance" desciption:"load_balance"`
	AccessControl *AccessControl `json:"access_control" desciption:"access_control"`
}

type ServicePercentage struct {
	ServiceType  int   `json:"service_type" description:"服务类型"`
	ServiceCount int64 `json:"service_count" description:"服务个数"`
}

type HttpServicePercentage struct {
	HttpServiceType  int   `json:"http_service_type" description:"HTTP服务类型"`
	HttpServiceCount int64 `json:"http_service_count" description:"HTTP服务个数"`
}
