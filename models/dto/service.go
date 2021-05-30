package dto

import (
	"github.com/LotteWong/giotto-gateway-core/utils"
	"github.com/gin-gonic/gin"
)

type ListServicesReq struct {
	Keyword   string `json:"keyword" form:"keyword" comment:"关键词" example:"keyword" validate:""` // 关键词
	PageIndex int    `json:"page_index" form:"page_index" comment:"当前页" example:"1" validate:""` // 当前页
	PageSize  int    `json:"page_size" form:"page_size" comment:"页条数" example:"1" validate:""`   // 页条数
}

type ListServicesRes struct {
	Total int64             `json:"total" form:"total" comment:"共计条数" example:"0" validate:""` // 共计条数
	Items []ListServiceItem `json:"items" form:"items" comment:"服务列表" example:"" validate:""`  // 服务列表
}

type ListServiceItem struct {
	Id          int64  `json:"id" form:"id"`
	ServiceName string `json:"service_name" form:"service_name"`
	ServiceDesc string `json:"service_desc" form:"service_desc"`
	ServiceType int    `json:"service_type" form:"service_type"`
	ServiceAddr string `json:"service_addr" form:"service_addr"`
	RealQps     int64  `json:"real_qps" form:"real_qps"`
	RealQpd     int64  `json:"real_qpd" form:"real_qpd"`
	TotalNode   int    `json:"total_node" form:"total_node"`
}

func (params *ListServicesReq) BindAndValidListServicesReq(ctx *gin.Context) error {
	return utils.DefaultGetValidParams(ctx, params)
}

type CreateOrUpdateHttpServiceReq struct {
	ServiceName string `json:"service_name" form:"service_name" comment:"服务名称" example:"" validate:"required,valid_service_name"` // 服务名称
	ServiceDesc string `json:"service_desc" form:"service_desc" comment:"服务描述" example:"" validate:"max=255,min=0"`               // 服务描述

	RuleType        int    `json:"rule_type" form:"rule_type" comment:"接入类型" example:"" validate:"max=1,min=0"`                              // 接入类型
	Rule            string `json:"rule" form:"rule" comment:"域名或者前缀" example:"" validate:"required,valid_rule"`                              // 域名或者前缀
	NeedHttps       int    `json:"need_https" form:"need_https" comment:"支持https" example:"" validate:"max=1,min=0"`                         // 支持https
	NeedStripUri    int    `json:"need_strip_uri" form:"need_strip_uri" comment:"启用strip_uri" example:"" validate:"max=1,min=0"`             // 启用strip_uri
	NeedWebsocket   int    `json:"need_websocket" form:"need_websocket" comment:"支持websocket" example:"" validate:"max=1,min=0"`             // 支持websocket
	UrlRewrite      string `json:"url_rewrite" form:"url_rewrite" comment:"url重写功能" example:"" validate:"valid_url_rewrite"`                 // 启用url重写
	HeaderTransform string `json:"header_transform" form:"header_transform" comment:"header转换" example:"" validate:"valid_header_transform"` // header转换

	OpenAuth  int    `json:"open_auth" form:"open_auth" comment:"是否开启权限" example:"" validate:"max=1,min=0"` // 是否开启权限
	BlackList string `json:"black_list" form:"black_list" comment:"黑名单ip" example:"" validate:""`           // 黑名单ip列表
	WhiteList string `json:"white_list" form:"white_list" comment:"白名单ip" example:"" validate:""`           // 白名单ip列表

	ClientIpFlowLimit       int64 `json:"client_ip_flow_limit" form:"client_ip_flow_limit" comment:"客户端ip限流数量" example:"" validate:"min=0"`             // 客户端ip限流数量
	ClientIpFlowInterval    int   `json:"client_ip_flow_interval" form:"client_ip_flow_interval" comment:"客户端ip限流间隔	" example:"" validate:"min=0"`      // 客户端ip限流间隔
	ServiceHostFlowLimit    int64 `json:"service_host_flow_limit" form:"service_host_flow_limit" comment:"服务端主机限流" example:"" validate:"min=0"`         // 服务端主机限流数量
	ServiceHostFlowInterval int   `json:"service_host_flow_interval" form:"service_host_flow_interval" comment:"服务端主机限流间隔" example:"" validate:"min=0"` // 客户端ip限流间隔

	RoundType  int    `json:"round_type" form:"round_type" comment:"轮询方式" example:"" validate:"max=3,min=0"`            // 轮询方式
	IpList     string `json:"ip_list" form:"ip_list" comment:"ip列表" example:"" validate:"valid_ip_port_list"`           // 启用ip列表
	WeightList string `json:"weight_list" form:"weight_list" comment:"权重列表" example:"" validate:"valid_weight_list"`    // 权重列表
	ForbidList string `json:"forbid_list" form:"forbid_list" comment:"禁用ip列表" example:"" validate:"valid_ip_port_list"` // 禁用ip列表

	UpstreamConnectTimeout int `json:"upstream_connect_timeout" form:"upstream_connect_timeout" comment:"建立连接超时时间, 单位为s" example:"" validate:"min=0"`   // 建立连接超时时间, 单位为s
	UpstreamHeaderTimeout  int `json:"upstream_header_timeout" form:"upstream_header_timeout" comment:"获取header超时时间, 单位为s" example:"" validate:"min=0"` // 获取header超时时间, 单位为s
	UpstreamIdleTimeout    int `json:"upstream_idle_timeout" form:"upstream_idle_timeout" comment:"连接最大空闲时间, 单位为s" example:"" validate:"min=0"`         // 连接最大空闲时间, 单位为s
	UpstreamMaxIdle        int `json:"upstream_max_idle" form:"upstream_max_idle" comment:"最大连接空闲数量" example:"" validate:"min=0"`                       // 最大连接空闲数量
}

func (params *CreateOrUpdateHttpServiceReq) BindAndValidCreateOrUpdateHttpServiceReq(ctx *gin.Context) error {
	return utils.DefaultGetValidParams(ctx, params)
}

type CreateOrUpdateTcpServiceReq struct {
	ServiceName string `json:"service_name" form:"service_name" comment:"服务名称" validate:"required,valid_service_name"` // 服务名称
	ServiceDesc string `json:"service_desc" form:"service_desc" comment:"服务描述" validate:"max=255,min=0"`               // 服务描述
	Port        int    `json:"port" form:"port" comment:"端口，需要设置在8001-8999范围内" validate:"required,min=8001,max=8999"`  // 端口

	OpenAuth      int    `json:"open_auth" form:"open_auth" comment:"是否开启权限" example:"" validate:"max=1,min=0"`     // 是否开启权限
	BlackList     string `json:"black_list" form:"black_list" comment:"黑名单ip" example:"" validate:"valid_ip_list"`  // 黑名单ip列表
	WhiteList     string `json:"white_list" form:"white_list" comment:"白名单ip" example:"" validate:"valid_ip_list"`  // 白名单ip列表
	WhiteHostName string `json:"white_host_name" form:"white_host_name" comment:"白名单主机列表" validate:"valid_ip_list"` // 白名单主机列表
	BlackHostName string `json:"black_host_name" form:"black_host_name" comment:"黑名单主机列表" validate:"valid_ip_list"` // 黑名单主机列表

	ClientIpFlowLimit       int64 `json:"client_ip_flow_limit" form:"client_ip_flow_limit" comment:"客户端ip限流数量" example:"" validate:"min=0"`             // 客户端ip限流数量
	ClientIpFlowInterval    int   `json:"client_ip_flow_interval" form:"client_ip_flow_interval" comment:"客户端ip限流间隔	" example:"" validate:"min=0"`      // 客户端ip限流间隔
	ServiceHostFlowLimit    int64 `json:"service_host_flow_limit" form:"service_host_flow_limit" comment:"服务端主机限流" example:"" validate:"min=0"`         // 服务端主机限流数量
	ServiceHostFlowInterval int   `json:"service_host_flow_interval" form:"service_host_flow_interval" comment:"服务端主机限流间隔" example:"" validate:"min=0"` // 客户端ip限流间隔

	RoundType  int    `json:"round_type" form:"round_type" comment:"轮询方式" example:"" validate:"max=3,min=0"`         // 轮询方式
	IpList     string `json:"ip_list" form:"ip_list" comment:"ip列表" example:"" validate:"valid_ip_port_list"`        // 启用ip列表
	WeightList string `json:"weight_list" form:"weight_list" comment:"权重列表" example:"" validate:"valid_weight_list"` // 权重列表
	ForbidList string `json:"forbid_list" form:"forbid_list" comment:"禁用ip列表" validate:"valid_ip_list"`              // 禁用ip列表
}

func (params *CreateOrUpdateTcpServiceReq) BindAndValidCreateOrUpdateTcpServiceReq(ctx *gin.Context) error {
	return utils.DefaultGetValidParams(ctx, params)
}

type CreateOrUpdateGrpcServiceReq struct {
	ServiceName     string `json:"service_name" form:"service_name" comment:"服务名称" validate:"required,valid_service_name"`                   // 服务名称
	ServiceDesc     string `json:"service_desc" form:"service_desc" comment:"服务描述" validate:"max=255,min=0"`                                 // 服务描述
	Port            int    `json:"port" form:"port" comment:"端口，需要设置在8001-8999范围内" validate:"required,min=8001,max=8999"`                    // 端口
	HeaderTransform string `json:"header_transform" form:"header_transform" comment:"header转换" example:"" validate:"valid_header_transform"` // header转换

	OpenAuth      int    `json:"open_auth" form:"open_auth" comment:"是否开启权限" example:"" validate:"max=1,min=0"`     // 是否开启权限
	BlackList     string `json:"black_list" form:"black_list" comment:"黑名单ip" example:"" validate:"valid_ip_list"`  // 黑名单ip列表
	WhiteList     string `json:"white_list" form:"white_list" comment:"白名单ip" example:"" validate:"valid_ip_list"`  // 白名单ip列表
	WhiteHostName string `json:"white_host_name" form:"white_host_name" comment:"白名单主机列表" validate:"valid_ip_list"` // 白名单主机列表
	BlackHostName string `json:"black_host_name" form:"black_host_name" comment:"黑名单主机列表" validate:"valid_ip_list"` // 黑名单主机列表

	ClientIpFlowLimit       int64 `json:"client_ip_flow_limit" form:"client_ip_flow_limit" comment:"客户端ip限流数量" example:"" validate:"min=0"`             // 客户端ip限流数量
	ClientIpFlowInterval    int   `json:"client_ip_flow_interval" form:"client_ip_flow_interval" comment:"客户端ip限流间隔	" example:"" validate:"min=0"`      // 客户端ip限流间隔
	ServiceHostFlowLimit    int64 `json:"service_host_flow_limit" form:"service_host_flow_limit" comment:"服务端主机限流" example:"" validate:"min=0"`         // 服务端主机限流数量
	ServiceHostFlowInterval int   `json:"service_host_flow_interval" form:"service_host_flow_interval" comment:"服务端主机限流间隔" example:"" validate:"min=0"` // 客户端ip限流间隔

	RoundType  int    `json:"round_type" form:"round_type" comment:"轮询方式" example:"" validate:"max=3,min=0"`         // 轮询方式
	IpList     string `json:"ip_list" form:"ip_list" comment:"ip列表" example:"" validate:"valid_ip_port_list"`        // 启用ip列表
	WeightList string `json:"weight_list" form:"weight_list" comment:"权重列表" example:"" validate:"valid_weight_list"` // 权重列表
	ForbidList string `json:"forbid_list" form:"forbid_list" comment:"禁用ip列表" validate:"valid_ip_list"`              // 禁用ip列表
}

func (params *CreateOrUpdateGrpcServiceReq) BindAndValidCreateOrUpdateGrpcServiceReq(ctx *gin.Context) error {
	return utils.DefaultGetValidParams(ctx, params)
}
