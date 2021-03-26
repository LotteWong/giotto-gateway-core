package dto

import (
	"github.com/LotteWong/giotto-gateway/utils"
	"github.com/gin-gonic/gin"
)

type JwtReq struct {
	Type       string `json:"type" form:"type" comment:"令牌类型" example:"" validate:""`             // 令牌类型
	Permission string `json:"permission" form:"permission" comment:"读写权限" example:"" validate:""` // 读写权限
}

type JwtRes struct {
	Token      string `json:"token" form:"token" comment:"令牌内容" example:"" validate:""`           // 令牌内容
	Type       string `json:"type" form:"type" comment:"令牌类型" example:"" validate:""`             // 令牌类型
	Permission string `json:"permission" form:"permission" comment:"读写权限" example:"" validate:""` // 读写权限
	ExpireAt   int    `json:"expire_at" form:"expire_at" comment:"失效时间" example:"" validate:""`   // 失效时间
}

func (params *JwtReq) BindAndValidJwtReq(ctx *gin.Context) error {
	return utils.DefaultGetValidParams(ctx, params)
}
