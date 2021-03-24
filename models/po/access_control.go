package po

type AccessControl struct {
	Id            int64  `json:"id" gorm:"primary_key" description:"自增主键"`
	ServiceId     int64  `json:"service_id" gorm:"column:service_id" description:"服务id"`
	OpenAuth      int    `json:"open_auth" gorm:"column:open_auth" description:"是否开启权限 1=开启"`
	BlackList     string `json:"black_list" gorm:"column:black_list" description:"黑名单ip"`
	WhiteList     string `json:"white_list" gorm:"column:white_list" description:"白名单ip"`
	WhiteHostName string `json:"white_host_name" gorm:"column:white_host_name" description:"白名单主机"`
	//BlackHostName     string `json:"black_host_name" gorm:"column:black_host_name" description:"黑名单主机"`
	ClientIpFlowLimit    int64 `json:"client_ip_flow_limit" gorm:"column:client_ip_flow_limit" description:"客户端ip限流"`
	ServiceHostFlowLimit int64 `json:"service_host_flow_limit" gorm:"column:service_host_flow_limit" description:"服务端主机限流"`
}

func (t *AccessControl) TableName() string {
	return "gateway_service_access_control"
}
