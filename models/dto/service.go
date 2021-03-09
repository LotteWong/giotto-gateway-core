package dto

import (
	"github.com/LotteWong/giotto-gateway/utils"
	"github.com/gin-gonic/gin"
)

type ListServicesReq struct {
	Keyword   string `json:"keyword" form:"keyword" comment:"关键词" example:"keyword" validate:""` // 关键词
	PageIndex int    `json:"page_index" form:"page_index" comment:"当前页" example:"1" validate:""` // 当前页
	PageSize  int    `json:"page_size" form:"page_size" comment:"页条数" example:"1" validate:""`   // 页条数
}

type ListServicesRes struct {
	Total int64     `json:"total" form:"total" comment:"共计条数" example:"0" validate:""` // 共计条数
	Items []Service `json:"items" form:"items" comment:"服务列表" example:"" validate:""`  // 服务列表
}

type Service struct {
	Id          int64  `json:"id" form:"id"`
	ServiceName string `json:"service_name" form:"service_name"`
	ServiceDesc string `json:"service_desc" form:"service_desc"`
	ServiceType int    `json:"service_type" form:"service_type"`
	ServiceAddr string `json:"service_addr" form:"service_addr"`
	Qps         int64  `json:"qps" form:"qps"`
	Qpd         int64  `json:"qpd" form:"qpd"`
	TotalNode   int    `json:"total_node" form:"total_node"`
}

func (params *ListServicesReq) BindAndValidListServicesReq(ctx *gin.Context) error {
	return utils.DefaultGetValidParams(ctx, params)
}
