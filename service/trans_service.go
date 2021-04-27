package service

import (
	"crypto/tls"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/LotteWong/giotto-gateway-core/constants"
	"github.com/LotteWong/giotto-gateway-core/models/po"
)

var transService *TransService

type TransService struct {
	TransportMap   map[string]*po.TransportDetail
	TransportSlice []*po.TransportDetail
	RWLock         sync.RWMutex
}

func NewTransService() *TransService {
	service := &TransService{
		TransportMap:   map[string]*po.TransportDetail{},
		TransportSlice: []*po.TransportDetail{},
		RWLock:         sync.RWMutex{},
	}
	return service
}

func GetTransService() *TransService {
	if transService == nil {
		transService = NewTransService()
	}
	return transService
}

func (s *TransService) GetTransForSvc(svc *po.ServiceDetail) (*http.Transport, error) {
	// hit in cache, use cache data
	for _, trans := range s.TransportSlice {
		if trans.ServiceName == svc.Info.ServiceName {
			return trans.Transporter, nil
		}
	}

	// miss in cache, new a transport
	if svc.LoadBalance.UpstreamConnectTimeout == 0 {
		svc.LoadBalance.UpstreamConnectTimeout = constants.DefaultUpstreamConnectTimeout
	}
	if svc.LoadBalance.UpstreamHeaderTimeout == 0 {
		svc.LoadBalance.UpstreamHeaderTimeout = constants.DefaultUpstreamHeaderTimeout
	}
	if svc.LoadBalance.UpstreamIdleTimeout == 0 {
		svc.LoadBalance.UpstreamIdleTimeout = constants.DefaultUpstreamIdleTimeout
	}
	if svc.LoadBalance.UpstreamMaxIdle == 0 {
		svc.LoadBalance.UpstreamMaxIdle = constants.DefaultUpstreamMaxIdle
	}

	transr := &http.Transport{
		Proxy:             http.ProxyFromEnvironment,
		ForceAttemptHTTP2: true, // support http2 protocol
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(svc.LoadBalance.UpstreamConnectTimeout) * time.Second,
			KeepAlive: time.Duration(constants.DefaultKeepAlive) * time.Second,
			DualStack: true,
		}).DialContext,
		ResponseHeaderTimeout: time.Duration(svc.LoadBalance.UpstreamHeaderTimeout) * time.Second,
		MaxIdleConns:          svc.LoadBalance.UpstreamMaxIdle,
		IdleConnTimeout:       time.Duration(svc.LoadBalance.UpstreamIdleTimeout) * time.Second,
		TLSHandshakeTimeout:   time.Duration(constants.DefaultTLSHandshakeTimeout) * time.Second,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}

	// miss in cache, write back to cache
	trans := &po.TransportDetail{
		Transporter: transr,
		ServiceName: svc.Info.ServiceName,
	}

	s.TransportSlice = append(s.TransportSlice, trans)
	s.RWLock.Lock()
	defer s.RWLock.Unlock()
	s.TransportMap[svc.Info.ServiceName] = trans

	return transr, nil
}
