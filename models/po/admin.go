package po

import "time"

type Admin struct {
	Id        int       `json:"id" gorm:"primary_key" description:"自增主键"`
	Username  string    `json:"username" gorm:"column:username" description:"登录名称"`
	Password  string    `json:"password" gorm:"column:password" description:"登录密码"`
	Salt      string    `json:"salt" gorm:"column:salt" description:"加密盐值"`
	UpdatedAt time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	CreatedAt time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	IsDelete  int       `json:"is_delete" gorm:"column:is_delete" description:"是否删除"`
}

func (t *Admin) TableName() string {
	return "gateway_admin"
}
