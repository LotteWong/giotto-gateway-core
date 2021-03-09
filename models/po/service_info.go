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
