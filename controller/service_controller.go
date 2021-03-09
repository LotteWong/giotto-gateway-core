package controller

import (
	"github.com/LotteWong/giotto-gateway/middleware"
	"github.com/LotteWong/giotto-gateway/models/dto"
	"github.com/LotteWong/giotto-gateway/service"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ServiceController struct{}

func RegistServiceRoutes(grp *gin.RouterGroup) {
	controller := &ServiceController{}
	grp.GET("", controller.ListServices)
	grp.DELETE("/:service_id", controller.DeleteService)
}

// ListServices godoc
// @Summary 查询服务列表接口
// @Description 查询服务列表
// @Tags 服务接口
// @Id /services
// @Produce  json
// @Param keyword query string false "keyword"
// @Param page_index query string false "page index"
// @Param page_size query string false "page size"
// @Success 200 {object} middleware.Response{data=dto.ListServicesRes} "success"
// @Router /services [get]
func (c *ServiceController) ListServices(ctx *gin.Context) {
	// validate request params
	req := &dto.ListServicesReq{}
	if err := req.BindAndValidListServicesReq(ctx); err != nil {
		middleware.ResponseError(ctx, 3000, err)
		return
	}

	// fuzzy search and page services
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(ctx, 3001, err)
		return
	}
	total, items, err := service.GetSvcService().ListServices(ctx, tx, req)
	if err != nil {
		middleware.ResponseError(ctx, 3002, err)
		return
	}

	// return response body
	res := &dto.ListServicesRes{
		Total: total,
		Items: items,
	}
	middleware.ResponseSuccess(ctx, res)
}

// DeleteService godoc
// @Summary 删除服务接口
// @Description 删除服务
// @Tags 服务接口
// @Id /services/{service_id}
// @Produce  json
// @Param service_id path string true "keyword"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /services/{service_id} [delete]
func (c *ServiceController) DeleteService(ctx *gin.Context) {
	serviceId, err := strconv.Atoi(ctx.Param("service_id"))
	if err != nil {
		middleware.ResponseError(ctx, 3000, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(ctx, 3001, err)
		return
	}
	err = service.GetSvcService().DeleteServices(ctx, tx, int64(serviceId))
	if err != nil {
		middleware.ResponseError(ctx, 3002, err)
		return
	}

	middleware.ResponseSuccess(ctx, "")
}
