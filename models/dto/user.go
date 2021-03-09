package dto

import (
	"github.com/LotteWong/giotto-gateway/utils"
	"github.com/gin-gonic/gin"
	"time"
)

type UserInfo struct {
	ID       int       `json:"id"`
	Username string    `json:"username"`
	LoginAt  time.Time `json:"login_at"`
	Avatar   string    `json:"avatar"`
	Intro    string    `json:"intro"`
	Roles    []string  `json:"roles"`
}

type ChangeUserPwdReq struct {
	Password string `json:"password" form:"password" comment:"登录密码" example:"123456" validate:"required"`
} // 登录密码

type ChangeUserPwdRes struct {
}

func (params *ChangeUserPwdReq) BindAndValidChangeUserPwdReq(ctx *gin.Context) error {
	return utils.DefaultGetValidParams(ctx, params)
}
