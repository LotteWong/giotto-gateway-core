package controller

import (
	"github.com/LotteWong/giotto-gateway/middleware"
	"github.com/LotteWong/giotto-gateway/models/dto"
	"github.com/LotteWong/giotto-gateway/service"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
)

type LoginController struct{}

func RegistLoginRoutes(grp *gin.RouterGroup) {
	controller := &LoginController{}
	grp.POST("/login", controller.Login)
	grp.POST("/logout", controller.Logout)
}

// Login godoc
// @Summary 用户登录接口
// @Description 使用登录名称和登录密码来登录
// @Tags 登录接口
// @Id /login
// @Accept  json
// @Produce  json
// @Param body body dto.LoginReq true "login request body"
// @Success 200 {object} middleware.Response{data=dto.LoginRes} "success"
// @Router /login [post]
func (c *LoginController) Login(ctx *gin.Context) {
	// validate request params
	req := &dto.LoginReq{}
	if err := req.BindAndValidLoginReq(ctx); err != nil {
		middleware.ResponseError(ctx, 1000, err)
		return
	}

	// login business logic
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(ctx, 1001, err)
		return
	}
	user, err := service.GetLoginService().Login(ctx, tx, req)
	if err != nil {
		middleware.ResponseError(ctx, 1002, err)
		return
	}

	// return response body
	res := &dto.LoginRes{Token: user.Username}
	middleware.ResponseSuccess(ctx, res)
}

// Logout godoc
// @Summary 用户登出接口
// @Description 退出登录并且删除会话信息
// @Tags 登录接口
// @Id /logout
// @Produce  json
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /logout [post]
func (c *LoginController) Logout(ctx *gin.Context) {
	if err := service.GetLoginService().Logout(ctx); err != nil {
		middleware.ResponseError(ctx, 1000, err)
		return
	}

	middleware.ResponseSuccess(ctx, "")
}
