package controller

import (
	"github.com/LotteWong/giotto-gateway/middleware"
	"github.com/LotteWong/giotto-gateway/models/dto"
	"github.com/LotteWong/giotto-gateway/service"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type ServiceController struct{}

func RegistServiceRoutes(grp *gin.RouterGroup) {
	controller := &ServiceController{}
	grp.GET("", controller.ListServices)
	grp.GET("/:service_id", controller.ShowService)
	grp.GET("/:service_id/status", controller.GetServiceStatus)
	grp.POST("/http", controller.CreateHttpService)
	grp.POST("/tcp", controller.CreateTcpService)
	grp.PUT("/tcp/:service_id", controller.UpdateTcpService)
	grp.POST("/grpc", controller.CreateGrpcService)
	grp.PUT("/grpc/:service_id", controller.UpdateGrpcService)
	grp.PUT("/http/:service_id", controller.UpdateHttpService)
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

// ShowService godoc
// @Summary 查询http服务详情接口
// @Description 查询http服务详情
// @Tags 服务接口
// @Id /services/{service_id}
// @Produce  json
// @Param service_id path string true "service id"
// @Success 200 {object} middleware.Response{data=po.ServiceDetail} "success"
// @Router /services/{service_id} [get]
func (c *ServiceController) ShowService(ctx *gin.Context) {
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
	res, err := service.GetSvcService().ShowService(ctx, tx, int64(serviceId))
	if err != nil {
		middleware.ResponseError(ctx, 3002, err)
		return
	}

	middleware.ResponseSuccess(ctx, res)
}

// GetServiceStatus godoc
// @Summary 查询http服务状态接口
// @Description 查询http服务状态
// @Tags 服务接口
// @Id /services/{service_id}/status
// @Produce  json
// @Param service_id path string true "service id"
// @Success 200 {object} middleware.Response{data=dto.ServiceStatus} "success"
// @Router /services/{service_id}/status [get]
func (c *ServiceController) GetServiceStatus(ctx *gin.Context) {
	_, err := strconv.Atoi(ctx.Param("service_id"))
	if err != nil {
		middleware.ResponseError(ctx, 3000, err)
		return
	}

	//tx, err := lib.GetGormPool("default")
	//if err != nil {
	//	middleware.ResponseError(ctx, 3001, err)
	//	return
	//}
	var todayFlow []int64
	var yesterdayFlow []int64
	for i := 0; i <= time.Now().Hour(); i++ {
		todayFlow = append(todayFlow, 0)
	}
	for i := 0; i <= 23; i++ {
		yesterdayFlow = append(yesterdayFlow, 0)
	}

	res := &dto.ServiceStatus{
		TodayFlow:     todayFlow,
		YesterdayFlow: yesterdayFlow,
	}
	middleware.ResponseSuccess(ctx, res)
}

// CreateHttpService godoc
// @Summary 创建http服务接口
// @Description 创建http服务
// @Tags 服务接口
// @Id /services/http
// @Accept  json
// @Produce  json
// @Param body body dto.CreateOrUpdateHttpServiceReq true "create http service request body"
// @Success 200 {object} middleware.Response{data=po.ServiceDetail} "success"
// @Router /services/http [post]
func (c *ServiceController) CreateHttpService(ctx *gin.Context) {
	// validate request params
	req := &dto.CreateOrUpdateHttpServiceReq{}
	if err := req.BindAndValidCreateOrUpdateHttpServiceReq(ctx); err != nil {
		middleware.ResponseError(ctx, 3000, err)
		return
	}

	// create http service
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(ctx, 3001, err)
		return
	}
	res, err := service.GetSvcService().CreateHttpService(ctx, tx, req)
	if err != nil {
		middleware.ResponseError(ctx, 3002, err)
		return
	}

	// return response body
	middleware.ResponseSuccess(ctx, res)
}

// UpdateHttpService godoc
// @Summary 更新http服务接口
// @Description 更新http服务
// @Tags 服务接口
// @Id /services/http/{service_id}
// @Accept  json
// @Produce  json
// @Param service_id path string true "service id"
// @Param body body dto.CreateOrUpdateHttpServiceReq true "update http service request body"
// @Success 200 {object} middleware.Response{data=po.ServiceDetail} "success"
// @Router /services/http/{service_id} [put]
func (c *ServiceController) UpdateHttpService(ctx *gin.Context) {
	// validate request params
	serviceId, err := strconv.Atoi(ctx.Param("service_id"))
	if err != nil {
		middleware.ResponseError(ctx, 3000, err)
		return
	}
	req := &dto.CreateOrUpdateHttpServiceReq{}
	if err := req.BindAndValidCreateOrUpdateHttpServiceReq(ctx); err != nil {
		middleware.ResponseError(ctx, 3001, err)
		return
	}

	// update http service
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(ctx, 3002, err)
		return
	}
	res, err := service.GetSvcService().UpdateHttpService(ctx, tx, req, int64(serviceId))
	if err != nil {
		middleware.ResponseError(ctx, 3003, err)
		return
	}

	// return response body
	middleware.ResponseSuccess(ctx, res)
}

// CreateTcpService godoc
// @Summary 创建tcp服务接口
// @Description 创建tcp服务
// @Tags 服务接口
// @Id /services/tcp
// @Accept  json
// @Produce  json
// @Param body body dto.CreateOrUpdateTcpServiceReq true "create tcp service request body"
// @Success 200 {object} middleware.Response{data=po.ServiceDetail} "success"
// @Router /services/tcp [post]
func (c *ServiceController) CreateTcpService(ctx *gin.Context) {
	// validate request params
	req := &dto.CreateOrUpdateTcpServiceReq{}
	if err := req.BindAndValidCreateOrUpdateTcpServiceReq(ctx); err != nil {
		middleware.ResponseError(ctx, 3000, err)
		return
	}

	// create http service
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(ctx, 3001, err)
		return
	}
	res, err := service.GetSvcService().CreateTcpService(ctx, tx, req)
	if err != nil {
		middleware.ResponseError(ctx, 3002, err)
		return
	}

	// return response body
	middleware.ResponseSuccess(ctx, res)
}

// UpdateTcpService godoc
// @Summary 更新tcp服务接口
// @Description 更新tcp服务
// @Tags 服务接口
// @Id /services/tcp/{service_id}
// @Accept  json
// @Produce  json
// @Param service_id path string true "service id"
// @Param body body dto.CreateOrUpdateTcpServiceReq true "update tcp service request body"
// @Success 200 {object} middleware.Response{data=po.ServiceDetail} "success"
// @Router /services/tcp/{service_id} [put]
func (c *ServiceController) UpdateTcpService(ctx *gin.Context) {
	// validate request params
	serviceId, err := strconv.Atoi(ctx.Param("service_id"))
	if err != nil {
		middleware.ResponseError(ctx, 3000, err)
		return
	}
	req := &dto.CreateOrUpdateTcpServiceReq{}
	if err := req.BindAndValidCreateOrUpdateTcpServiceReq(ctx); err != nil {
		middleware.ResponseError(ctx, 3001, err)
		return
	}

	// update http service
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(ctx, 3002, err)
		return
	}
	res, err := service.GetSvcService().UpdateTcpService(ctx, tx, req, int64(serviceId))
	if err != nil {
		middleware.ResponseError(ctx, 3003, err)
		return
	}

	// return response body
	middleware.ResponseSuccess(ctx, res)
}

// CreateGrpcService godoc
// @Summary 创建grpc服务接口
// @Description 创建grpc服务
// @Tags 服务接口
// @Id /services/grpc
// @Accept  json
// @Produce  json
// @Param body body dto.CreateOrUpdateGrpcServiceReq true "create grpc service request body"
// @Success 200 {object} middleware.Response{data=po.ServiceDetail} "success"
// @Router /services/grpc [post]
func (c *ServiceController) CreateGrpcService(ctx *gin.Context) {
	// validate request params
	req := &dto.CreateOrUpdateGrpcServiceReq{}
	if err := req.BindAndValidCreateOrUpdateGrpcServiceReq(ctx); err != nil {
		middleware.ResponseError(ctx, 3000, err)
		return
	}

	// create http service
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(ctx, 3001, err)
		return
	}
	res, err := service.GetSvcService().CreateGrpcService(ctx, tx, req)
	if err != nil {
		middleware.ResponseError(ctx, 3002, err)
		return
	}

	// return response body
	middleware.ResponseSuccess(ctx, res)
}

// UpdateGrpcService godoc
// @Summary 更新grpc服务接口
// @Description 更新grpc服务
// @Tags 服务接口
// @Id /services/grpc/{service_id}
// @Accept  json
// @Produce  json
// @Param service_id path string true "service id"
// @Param body body dto.CreateOrUpdateGrpcServiceReq true "update grpc service request body"
// @Success 200 {object} middleware.Response{data=po.ServiceDetail} "success"
// @Router /services/grpc/{service_id} [put]
func (c *ServiceController) UpdateGrpcService(ctx *gin.Context) {
	// validate request params
	serviceId, err := strconv.Atoi(ctx.Param("service_id"))
	if err != nil {
		middleware.ResponseError(ctx, 3000, err)
		return
	}
	req := &dto.CreateOrUpdateGrpcServiceReq{}
	if err := req.BindAndValidCreateOrUpdateGrpcServiceReq(ctx); err != nil {
		middleware.ResponseError(ctx, 3001, err)
		return
	}

	// update http service
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(ctx, 3002, err)
		return
	}
	res, err := service.GetSvcService().UpdateGrpcService(ctx, tx, req, int64(serviceId))
	if err != nil {
		middleware.ResponseError(ctx, 3003, err)
		return
	}

	// return response body
	middleware.ResponseSuccess(ctx, res)
}

// DeleteService godoc
// @Summary 删除服务接口
// @Description 删除服务
// @Tags 服务接口
// @Id /services/{service_id}
// @Produce  json
// @Param service_id path string true "service id"
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
