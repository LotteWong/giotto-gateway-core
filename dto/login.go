package dto

import (
	"github.com/LotteWong/giotto-gateway/utils"
	"github.com/gin-gonic/gin"
	"time"
)

type LoginReq struct {
	Username string `json:"username" form:"username" comment:"登录名称" example:"admin" validate:"required"`
	Password string `json:"password" form:"password" comment:"登录密码" example:"123456" validate:"required"`
}

type LoginRes struct {
	Token string `json:"token" form:"token" comment:"登录令牌" example:"token" validate:""`
}

type LoginSession struct {
	Id       int       `json:"id"`
	Username string    `json:"username"`
	LoginAt  time.Time `json:"login_at"`
}

func (params *LoginReq) BindAndValidLoginReq(ctx *gin.Context) error {
	return utils.DefaultGetValidParams(ctx, params)
}
