package controller

import (
	"github.com/LotteWong/giotto-gateway/constants"
	"github.com/LotteWong/giotto-gateway/middleware"
	"github.com/LotteWong/giotto-gateway/models/dto"
	"github.com/LotteWong/giotto-gateway/service"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type DashboardController struct{}

func RegistDashboardRoutes(grp *gin.RouterGroup) {
	controller := &DashboardController{}
	grp.GET("/flow", controller.GetTotalFlow)
	grp.GET("/flow/services/:service_id", controller.GetServiceFlow)
	grp.GET("/flow/apps/:app_id", controller.GetAppFlow)
	grp.GET("/statistics", controller.GetStatistics)
	grp.GET("/percentage/services", controller.GetServicePercentage)
}

// GetStatistics godoc
// @Summary 查询统计指标接口
// @Description 查询统计指标
// @Tags 数据接口
// @Id /statistics
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.Statistics} "success"
// @Router /statistics [get]
func (c *DashboardController) GetStatistics(ctx *gin.Context) {
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(ctx, 5000, err)
		return
	}

	serviceCount, _, err := service.GetSvcService().ListServices(ctx, tx, &dto.ListServicesReq{})
	if err != nil {
		middleware.ResponseError(ctx, 5001, err)
		return
	}

	appCount, _, err := service.GetAppService().ListApps(ctx, tx, &dto.ListAppsReq{})
	if err != nil {
		middleware.ResponseError(ctx, 5002, err)
		return
	}

	flowCount, err := service.GetFlowCountService().GetFlowCount(constants.TotalFlowCountPrefix)
	if err != nil {
		middleware.ResponseError(ctx, 5003, err)
		return
	}

	res := &dto.Statistics{
		ServiceCount: serviceCount,
		AppCount:     appCount,
		CurrentQpd:   flowCount.TotalCount,
		CurrentQps:   flowCount.Qps,
	}
	middleware.ResponseSuccess(ctx, res)
}

// GetServicePercentage godoc
// @Summary 查询服务类型占比接口
// @Description 查询服务类型占比
// @Tags 数据接口
// @Id /percentage/services
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.Flow} "success"
// @Router /percentage/services [get]
func (c *DashboardController) GetServicePercentage(ctx *gin.Context) {
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(ctx, 5000, err)
		return
	}

	res, err := service.GetDashboardService().GetServicePercentage(ctx, tx)
	if err != nil {
		middleware.ResponseError(ctx, 5001, err)
		return
	}

	middleware.ResponseSuccess(ctx, res)
}

// GetTotalFlow godoc
// @Summary 查询全部流量接口
// @Description 查询全部流量
// @Tags 数据接口
// @Id /flow
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.Flow} "success"
// @Router /flow [get]
func (c *DashboardController) GetTotalFlow(ctx *gin.Context) {
	count, err := service.GetFlowCountService().GetFlowCount(constants.TotalFlowCountPrefix)
	if err != nil {
		middleware.ResponseError(ctx, 5000, err)
		return
	}

	var todayFlow []int64
	todayTime := time.Now()
	for i := 0; i <= todayTime.Hour(); i++ {
		dateTime := time.Date(todayTime.Year(), todayTime.Month(), todayTime.Day(), i, 0, 0, 0, lib.TimeLocation)
		hourFlow, _ := service.GetFlowCountService().GetHourFlow(dateTime, count.ServiceName)
		todayFlow = append(todayFlow, hourFlow)
	}

	var yesterdayFlow []int64
	yeasterdayTime := todayTime.Add(-1 * 24 * time.Hour)
	for i := 0; i <= 23; i++ {
		dateTime := time.Date(yeasterdayTime.Year(), yeasterdayTime.Month(), yeasterdayTime.Day(), i, 0, 0, 0, lib.TimeLocation)
		hourFlow, _ := service.GetFlowCountService().GetHourFlow(dateTime, count.ServiceName)
		yesterdayFlow = append(yesterdayFlow, hourFlow)
	}

	res := &dto.Flow{
		TodayFlow:     todayFlow,
		YesterdayFlow: yesterdayFlow,
	}
	middleware.ResponseSuccess(ctx, res)
}

// GetServiceFlow godoc
// @Summary 查询服务流量接口
// @Description 查询服务流量
// @Tags 数据接口
// @Id /flow/services/{service_id}
// @Produce  json
// @Param service_id path string true "service id"
// @Success 200 {object} middleware.Response{data=dto.Flow} "success"
// @Router /flow/services/{service_id} [get]
func (c *DashboardController) GetServiceFlow(ctx *gin.Context) {
	service_id, err := strconv.Atoi(ctx.Param("service_id"))
	if err != nil {
		middleware.ResponseError(ctx, 5000, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(ctx, 5001, err)
		return
	}
	serviceDetail, err := service.GetSvcService().ShowService(ctx, tx, int64(service_id))
	if err != nil {
		middleware.ResponseError(ctx, 5002, err)
		return
	}

	count, err := service.GetFlowCountService().GetFlowCount(constants.ServiceFlowCountPrefix + serviceDetail.Info.ServiceName)
	if err != nil {
		middleware.ResponseError(ctx, 5003, err)
		return
	}

	var todayFlow []int64
	todayTime := time.Now()
	for i := 0; i <= todayTime.Hour(); i++ {
		dateTime := time.Date(todayTime.Year(), todayTime.Month(), todayTime.Day(), i, 0, 0, 0, lib.TimeLocation)
		hourFlow, _ := service.GetFlowCountService().GetHourFlow(dateTime, count.ServiceName)
		todayFlow = append(todayFlow, hourFlow)
	}

	var yesterdayFlow []int64
	yesterdayTime := todayTime.Add(-1 * 24 * time.Hour)
	for i := 0; i <= 23; i++ {
		dateTime := time.Date(yesterdayTime.Year(), yesterdayTime.Month(), yesterdayTime.Day(), i, 0, 0, 0, lib.TimeLocation)
		hourFlow, _ := service.GetFlowCountService().GetHourFlow(dateTime, count.ServiceName)
		yesterdayFlow = append(yesterdayFlow, hourFlow)
	}

	res := &dto.Flow{
		TodayFlow:     todayFlow,
		YesterdayFlow: yesterdayFlow,
	}
	middleware.ResponseSuccess(ctx, res)
}

// GetAppFlow godoc
// @Summary 查询租户流量接口
// @Description 查询租户流量
// @Tags 数据接口
// @Id /flow/apps/{app_id}
// @Produce  json
// @Param app_id path string true "app id"
// @Success 200 {object} middleware.Response{data=dto.Flow} "success"
// @Router /flow/apps/{app_id} [get]
func (c *DashboardController) GetAppFlow(ctx *gin.Context) {
	appId, err := strconv.Atoi(ctx.Param("app_id"))
	if err != nil {
		middleware.ResponseError(ctx, 5000, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(ctx, 5001, err)
		return
	}
	app, err := service.GetAppService().ShowApp(ctx, tx, int64(appId))
	if err != nil {
		middleware.ResponseError(ctx, 5002, err)
		return
	}

	count, err := service.GetFlowCountService().GetFlowCount(constants.AppFlowCountPrefix + app.AppId)
	if err != nil {
		middleware.ResponseError(ctx, 5003, err)
		return
	}

	var todayFlow []int64
	todayTime := time.Now()
	for i := 0; i <= todayTime.Hour(); i++ {
		dateTime := time.Date(todayTime.Year(), todayTime.Month(), todayTime.Day(), i, 0, 0, 0, lib.TimeLocation)
		hourFlow, _ := service.GetFlowCountService().GetHourFlow(dateTime, count.ServiceName)
		todayFlow = append(todayFlow, hourFlow)
	}

	var yesterdayFlow []int64
	yesterdayTime := todayTime.Add(-1 * 24 * time.Hour)
	for i := 0; i <= 23; i++ {
		dateTime := time.Date(yesterdayTime.Year(), yesterdayTime.Month(), yesterdayTime.Day(), i, 0, 0, 0, lib.TimeLocation)
		hourFlow, _ := service.GetFlowCountService().GetHourFlow(dateTime, count.ServiceName)
		yesterdayFlow = append(yesterdayFlow, hourFlow)
	}

	res := &dto.Flow{
		TodayFlow:     todayFlow,
		YesterdayFlow: yesterdayFlow,
	}
	middleware.ResponseSuccess(ctx, res)
}
