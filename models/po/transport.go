package po

import (
	"net/http"
)

type TransportDetail struct {
	Transporter *http.Transport `json:"transport" description:"下游传输器"`
	ServiceName string          `json:"service_name" description:"服务名称"`
}
