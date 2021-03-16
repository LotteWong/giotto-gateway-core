package po

import "time"

type App struct {
	Id        int64     `json:"id" gorm:"primary_key" description:"自增主键"`
	AppId     string    `json:"app_id" gorm:"column:app_id" description:"租户id"`
	AppName   string    `json:"app_name" gorm:"column:app_name" description:"租户名称"`
	Secret    string    `json:"secret" gorm:"column:secret" description:"租户密钥"`
	WhiteIps  string    `json:"white_ips" gorm:"column:white_ips" description:"ip白名单，支持前缀匹配"`
	Qpd       int64     `json:"qpd" gorm:"column:qpd" description:"每日请求量限制"`
	Qps       int64     `json:"qps" gorm:"column:qps" description:"每秒请求量限制"`
	CreatedAt time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	UpdatedAt time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	IsDelete  int       `json:"is_delete" gorm:"column:is_delete" description:"是否删除"`
}

func (t *App) TableName() string {
	return "gateway_app"
}
