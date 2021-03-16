package dto

type Flow struct {
	TodayFlow     []int64 `json:"today_flow" form:"today_flow" comment:"今日流量列表" example:"" validate:""`         // 今日流量列表
	YesterdayFlow []int64 `json:"yesterday_flow" form:"yesterday_flow" comment:"昨日流量列表" example:"" validate:""` // 昨日流量列表
}

type Statistics struct {
	ServiceCount int64 `json:"service_count" form:"service_count" comment:"服务总数" example:"" validate:""` // 服务总数
	AppCount     int64 `json:"app_count" form:"app_count" comment:"租户总数" example:"" validate:""`         // 租户总数
	CurrentQps   int64 `json:"current_qps" form:"current_qps" comment:"当前每秒请求量" example:"" validate:""`  // 当前每秒请求量
	CurrentQpd   int64 `json:"current_qpd" form:"current_qpd" comment:"当前每日请求量" example:"" validate:""`  // 当前每日请求量
}

type ServicePercentageItems struct {
	Legends []string                `json:"legends" form:"legends" comment:"图例列表" example:"" validate:""` // 图例列表
	Records []ServicePercentageItem `json:"records" form:"records" comment:"数据列表" example:"" validate:""` // 数据列表
}

type ServicePercentageItem struct {
	ServiceLegend string `json:"service_legend" form:"service_legend" comment:"服务图例" example:"" validate:""` // 服务图例
	ServiceType   int    `json:"service_type" form:"service_type" comment:"服务类型" example:"" validate:""`     // 服务类型
	ServiceCount  int64  `json:"service_count" form:"service_count" comment:"服务个数" example:"" validate:""`   // 服务个数
}
