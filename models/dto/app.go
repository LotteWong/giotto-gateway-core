package dto

import (
	"github.com/LotteWong/giotto-gateway/utils"
	"github.com/gin-gonic/gin"
)

type ListAppsReq struct {
	Keyword   string `json:"keyword" form:"keyword" comment:"关键词" example:"keyword" validate:""` // 关键词
	PageIndex int    `json:"page_index" form:"page_index" comment:"当前页" example:"1" validate:""` // 当前页
	PageSize  int    `json:"page_size" form:"page_size" comment:"页条数" example:"1" validate:""`   // 页条数
}

type ListAppsRes struct {
	Total int64         `json:"total" form:"total" comment:"共计条数" example:"0" validate:""` // 共计条数
	Items []ListAppItem `json:"items" form:"items" comment:"租户列表" example:"" validate:""`  // 租户列表
}

type ListAppItem struct {
	Id       int64  `json:"id" form:"id" comment:"自增主键" validate:""`
	AppId    string `json:"app_id" form:"app_id" comment:"租户id" validate:""`
	AppName  string `json:"name" form:"name" comment:"租户名称" validate:""`
	Secret   string `json:"secret" form:"secret" comment:"租户密钥" validate:""`
	WhiteIps string `json:"white_ips" form:"white_ips" comment:"ip白名单，支持前缀匹配" validate:""`
	Qpd      int64  `json:"qpd" form:"qpd" comment:"每日请求量限制" validate:""`
	Qps      int64  `json:"qps" form:"qps" comment:"每秒请求量限制" validate:""`
	RealQpd  int64  `json:"real_qpd" form:"real_qpd" comment:"每日实际请求量" validate:""`
	RealQps  int64  `json:"real_qps" form:"real_qps" comment:"每秒实际请求量" validate:""`
}

func (params *ListAppsReq) BindAndValidListAppsReq(ctx *gin.Context) error {
	return utils.DefaultGetValidParams(ctx, params)
}

type CreateOrUpdateAppReq struct {
	AppId    string `json:"app_id" form:"app_id" comment:"租户id" validate:""`               // 租户id
	AppName  string `json:"app_name" form:"app_name" comment:"租户名称" validate:""`           // 租户名称
	Secret   string `json:"secret" form:"secret" comment:"租户密钥" validate:""`               // 租户密钥
	WhiteIps string `json:"white_ips" form:"white_ips" comment:"ip白名单，支持前缀匹配" validate:""` // ip白名单，支持前缀匹配
	Qpd      int64  `json:"qpd" form:"qpd" comment:"每日请求量限制" validate:""`                  // 每日请求量限制
	Qps      int64  `json:"qps" form:"qps" comment:"每秒请求量限制" validate:""`                  // 每秒请求量限制
}

func (params *CreateOrUpdateAppReq) BindAndValidCreateOrUpdateAppReq(ctx *gin.Context) error {
	return utils.DefaultGetValidParams(ctx, params)
}
