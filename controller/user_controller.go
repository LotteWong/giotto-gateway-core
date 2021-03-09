package controller

import (
	"github.com/LotteWong/giotto-gateway/middleware"
	"github.com/LotteWong/giotto-gateway/models/dto"
	"github.com/LotteWong/giotto-gateway/service"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
)

type UserController struct{}

func RegistUserRoutes(grp *gin.RouterGroup) {
	controller := &UserController{}
	grp.GET("/admin", controller.GetUserInfo)
	grp.PATCH("/admin", controller.ChangeUserPwd)
}

// GetUserInfo godoc
// @Summary 查询用户信息接口
// @Description 获取用户的信息
// @Tags 用户接口
// @Id /users/admin
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.UserInfo} "success"
// @Router /users/admin [get]
func (c *UserController) GetUserInfo(ctx *gin.Context) {
	userInfo, err := service.GetUserService().GetUserInfo(ctx)
	if err != nil {
		middleware.ResponseError(ctx, 2000, err)
		return
	}

	middleware.ResponseSuccess(ctx, userInfo)
}

// ChangeUserPwd godoc
// @Summary 修改用户密码接口
// @Description 修改用户的密码
// @Tags 用户接口
// @Id /users/admin
// @Accept  json
// @Produce  json
// @Param body body dto.ChangeUserPwdReq true "change user password request body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /users/admin [POST]
func (c *UserController) ChangeUserPwd(ctx *gin.Context) {
	// validate request params
	req := &dto.ChangeUserPwdReq{}
	if err := req.BindAndValidChangeUserPwdReq(ctx); err != nil {
		middleware.ResponseError(ctx, 2000, err)
		return
	}

	// change user password business logic
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(ctx, 2001, err)
		return
	}
	if err := service.GetUserService().ChangeUserPassword(ctx, tx, req); err != nil {
		middleware.ResponseError(ctx, 2002, err)
		return
	}

	// return response body
	middleware.ResponseSuccess(ctx, "")
}
