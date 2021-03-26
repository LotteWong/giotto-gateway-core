package constants

const (
	ValidatorKey    = "ValidatorKey"
	TranslatorKey   = "TranslatorKey"
	LoginSessionKey = "LoginSessionKey"

	FlowHourCountKey = "flow_hour_count"
	FlowDayCountKey  = "flow_day_count"

	TotalFlowCountPrefix   = "flow_total_count"
	AppFlowCountPrefix     = "flow_app_count_"
	ServiceFlowCountPrefix = "flow_service_count_"

	JwtSignKey   = "jwt_sign_key"
	JwtExpires   = 60 * 60
	JwtType      = "Bearer"
	JwtReadWrite = "read-write"
	JwtReadOnly  = "read-only"
)
