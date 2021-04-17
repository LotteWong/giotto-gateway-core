package grpc_proxy_middleware

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/LotteWong/giotto-gateway/constants"
	"github.com/LotteWong/giotto-gateway/models/po"
	"github.com/LotteWong/giotto-gateway/service"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func GrpcRateLimitMiddleware(grpcServiceDetail *po.ServiceDetail) func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		metaCtx, ok := metadata.FromIncomingContext(stream.Context())
		if !ok {
			return errors.New("failed to get metadata context")
		}
		apps := metaCtx.Get("app")
		if len(apps) == 0 {
			return errors.New("app not found")
		}
		app := &po.App{}
		if err := json.Unmarshal([]byte(apps[0]), app); err != nil {
			return err
		}

		if grpcServiceDetail.AccessControl.ServiceHostFlowLimit != 0 {
			svrServiceName := constants.ServiceFlowCountPrefix + grpcServiceDetail.Info.ServiceName
			svrRateLimit, err := service.GetRateLimitService().GetRateLimit(svrServiceName, grpcServiceDetail.AccessControl.ServiceHostFlowLimit)
			if err != nil {
				return err
			}
			if !svrRateLimit.Allow() {
				return errors.New(fmt.Sprintf("service host flow limit is %d, rate limit exceeds", grpcServiceDetail.AccessControl.ServiceHostFlowLimit))
			}
		}

		if grpcServiceDetail.AccessControl.ClientIpFlowLimit != 0 {
			peerCtx, ok := peer.FromContext(stream.Context())
			if !ok {
				return errors.New("failed to get peer context")
			}
			peerAddr := peerCtx.Addr.String()
			clientIp := peerAddr[0:strings.LastIndex(peerAddr, ":")]

			cltServiceName := constants.ServiceFlowCountPrefix + grpcServiceDetail.Info.ServiceName + "_" + clientIp
			cltRateLimit, err := service.GetRateLimitService().GetRateLimit(cltServiceName, grpcServiceDetail.AccessControl.ClientIpFlowLimit)
			if err != nil {
				return err
			}
			if !cltRateLimit.Allow() {
				return errors.New(fmt.Sprintf("client ip flow limit is %d, rate limit exceeds", grpcServiceDetail.AccessControl.ClientIpFlowLimit))
			}
		}

		if app.Qps != 0 {
			appServiceName := constants.AppFlowCountPrefix + app.AppId
			appRateLimit, err := service.GetRateLimitService().GetRateLimit(appServiceName, app.Qps)
			if err != nil {
				return err
			}
			if !appRateLimit.Allow() {
				return errors.New(fmt.Sprintf("app flow limit is %d, rate limit exceeds", app.Qps))
			}
		}

		if err := handler(srv, stream); err != nil {
			return err
		}

		return nil
	}
}
