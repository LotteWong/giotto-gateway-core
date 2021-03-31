package controller

import (
	"encoding/base64"
	"fmt"
	"github.com/LotteWong/giotto-gateway/common_middleware"
	"github.com/LotteWong/giotto-gateway/models/dto"
	"github.com/LotteWong/giotto-gateway/service"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strings"
)

type JwtController struct{}

func RegistJwtRoutes(grp *gin.RouterGroup) {
	controller := &JwtController{}
	grp.POST("/jwt", controller.GenerateJwt)
}

// GenerateJwt godoc
// @Summary 生成认证令牌接口
// @Description 生成认证令牌
// @Tags 认证接口
// @Id /tokens/jwt
// @Accept  json
// @Produce  json
// @Param body body dto.JwtReq true "generate jwt request body"
// @Success 200 {object} common_middleware.Response{data=dto.JwtRes} "success"
// @Router /tokens/jwt [post]
func (c *JwtController) GenerateJwt(ctx *gin.Context) {
	// validate request params
	req := &dto.JwtReq{}
	if err := req.BindAndValidJwtReq(ctx); err != nil {
		common_middleware.ResponseError(ctx, 6000, err)
		return
	}

	// parse authorization to get app id and secret
	cipherInfo := strings.Split(ctx.GetHeader("Authorization"), " ")
	if len(cipherInfo) != 2 {
		common_middleware.ResponseError(ctx, 6001, errors.New("can not get app id and secret from authorization header"))
		return
	}

	plainInfo, err := base64.StdEncoding.DecodeString(cipherInfo[1])
	if err != nil {
		common_middleware.ResponseError(ctx, 6002, errors.New(fmt.Sprintf("base64 decode app id and secret error: %v", err)))
		return
	}

	pair := strings.Split(string(plainInfo), ":")
	if len(pair) != 2 {
		common_middleware.ResponseError(ctx, 6003, errors.New("can not get app id and secret from authorization header"))
		return
	}

	appId := pair[0]
	secret := pair[1]

	// generate jwt by app id and secret
	tx, err := lib.GetGormPool("default")
	if err != nil {
		common_middleware.ResponseError(ctx, 6004, err)
		return
	}
	res, err := service.GetJwtService().GenerateJwt(ctx, tx, req, appId, secret)
	if err != nil {
		common_middleware.ResponseError(ctx, 6005, err)
		return
	}

	// return response body
	common_middleware.ResponseSuccess(ctx, res)
}
