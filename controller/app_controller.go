package controller

import (
	"github.com/LotteWong/giotto-gateway/common_middleware"
	"github.com/LotteWong/giotto-gateway/models/dto"
	"github.com/LotteWong/giotto-gateway/service"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"strconv"
)

type AppController struct{}

func RegistAppRoutes(grp *gin.RouterGroup) {
	controller := &AppController{}
	grp.GET("", controller.ListApps)
	grp.GET("/:app_id", controller.ShowApp)
	grp.POST("", controller.CreateApp)
	grp.PUT("/:app_id", controller.UpdateApp)
	grp.DELETE("/:app_id", controller.DeleteApp)
}

// ListApps godoc
// @Summary 查询租户列表接口
// @Description 查询租户列表
// @Tags 租户接口
// @Id /apps
// @Produce  json
// @Param keyword query string false "keyword"
// @Param page_index query string false "page index"
// @Param page_size query string false "page size"
// @Success 200 {object} common_middleware.Response{data=dto.ListAppsRes} "success"
// @Router /apps [get]
func (c *AppController) ListApps(ctx *gin.Context) {
	// validate request params
	req := &dto.ListAppsReq{}
	if err := req.BindAndValidListAppsReq(ctx); err != nil {
		common_middleware.ResponseError(ctx, 4000, err)
		return
	}

	// fuzzy search and page apps
	tx, err := lib.GetGormPool("default")
	if err != nil {
		common_middleware.ResponseError(ctx, 4001, err)
		return
	}
	total, items, err := service.GetAppService().ListApps(ctx, tx, req)
	if err != nil {
		common_middleware.ResponseError(ctx, 4002, err)
		return
	}

	// return response body
	res := &dto.ListAppsRes{
		Total: total,
		Items: items,
	}
	common_middleware.ResponseSuccess(ctx, res)
}

// ShowApp godoc
// @Summary 查询租户详情接口
// @Description 查询租户详情
// @Tags 租户接口
// @Id /apps/{app_id}
// @Produce  json
// @Param app_id path string true "app id"
// @Success 200 {object} common_middleware.Response{data=po.App} "success"
// @Router /apps/{app_id} [get]
func (c *AppController) ShowApp(ctx *gin.Context) {
	appId, err := strconv.Atoi(ctx.Param("app_id"))
	if err != nil {
		common_middleware.ResponseError(ctx, 4000, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		common_middleware.ResponseError(ctx, 4001, err)
		return
	}
	res, err := service.GetAppService().ShowApp(ctx, tx, int64(appId))
	if err != nil {
		common_middleware.ResponseError(ctx, 4002, err)
		return
	}

	common_middleware.ResponseSuccess(ctx, res)
}

// CreateApp godoc
// @Summary 创建租户接口
// @Description 创建租户
// @Tags 租户接口
// @Id /apps
// @Accept  json
// @Produce  json
// @Param body body dto.CreateOrUpdateAppReq true "create app request body"
// @Success 200 {object} common_middleware.Response{data=po.App} "success"
// @Router /apps [post]
func (c *AppController) CreateApp(ctx *gin.Context) {
	// validate request params
	req := &dto.CreateOrUpdateAppReq{}
	if err := req.BindAndValidCreateOrUpdateAppReq(ctx); err != nil {
		common_middleware.ResponseError(ctx, 4000, err)
		return
	}

	// create http service
	tx, err := lib.GetGormPool("default")
	if err != nil {
		common_middleware.ResponseError(ctx, 4001, err)
		return
	}
	res, err := service.GetAppService().CreateApp(ctx, tx, req)
	if err != nil {
		common_middleware.ResponseError(ctx, 4002, err)
		return
	}

	// return response body
	common_middleware.ResponseSuccess(ctx, res)
}

// UpdateApp godoc
// @Summary 更新租户接口
// @Description 更新租户
// @Tags 租户接口
// @Id /apps/{app_id}
// @Accept  json
// @Produce  json
// @Param app_id path string true "app id"
// @Param body body dto.CreateOrUpdateAppReq true "update app request body"
// @Success 200 {object} common_middleware.Response{data=po.App} "success"
// @Router /apps/{app_id} [put]
func (c *AppController) UpdateApp(ctx *gin.Context) {
	// validate request params
	appId, err := strconv.Atoi(ctx.Param("app_id"))
	if err != nil {
		common_middleware.ResponseError(ctx, 4000, err)
		return
	}
	req := &dto.CreateOrUpdateAppReq{}
	if err := req.BindAndValidCreateOrUpdateAppReq(ctx); err != nil {
		common_middleware.ResponseError(ctx, 4001, err)
		return
	}

	// update http service
	tx, err := lib.GetGormPool("default")
	if err != nil {
		common_middleware.ResponseError(ctx, 4002, err)
		return
	}
	res, err := service.GetAppService().UpdateApp(ctx, tx, req, int64(appId))
	if err != nil {
		common_middleware.ResponseError(ctx, 4003, err)
		return
	}

	// return response body
	common_middleware.ResponseSuccess(ctx, res)
}

// DeleteService godoc
// @Summary 删除租户接口
// @Description 删除租户
// @Tags 租户接口
// @Id /apps/{app_id}
// @Produce  json
// @Param app_id path string true "app id"
// @Success 200 {object} common_middleware.Response{data=string} "success"
// @Router /apps/{app_id} [delete]
func (c *AppController) DeleteApp(ctx *gin.Context) {
	appId, err := strconv.Atoi(ctx.Param("app_id"))
	if err != nil {
		common_middleware.ResponseError(ctx, 4000, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		common_middleware.ResponseError(ctx, 4001, err)
		return
	}
	err = service.GetAppService().DeleteApp(ctx, tx, int64(appId))
	if err != nil {
		common_middleware.ResponseError(ctx, 4002, err)
		return
	}

	common_middleware.ResponseSuccess(ctx, "")
}
