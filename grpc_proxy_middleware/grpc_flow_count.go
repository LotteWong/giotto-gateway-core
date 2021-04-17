package grpc_proxy_middleware

import (
	"encoding/json"

	"github.com/LotteWong/giotto-gateway/constants"
	"github.com/LotteWong/giotto-gateway/models/po"
	"github.com/LotteWong/giotto-gateway/service"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func GrpcFlowCountMiddleware(grpcServiceDetail *po.ServiceDetail) func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
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

		ttlServiceName := constants.TotalFlowCountPrefix
		ttlFlowCount, err := service.GetFlowCountService().GetFlowCount(ttlServiceName)
		if err != nil {
			return err
		}
		service.GetFlowCountService().Increase(ttlFlowCount)

		svcServiceName := constants.ServiceFlowCountPrefix + grpcServiceDetail.Info.ServiceName
		svcFlowCount, err := service.GetFlowCountService().GetFlowCount(svcServiceName)
		if err != nil {
			return err
		}
		service.GetFlowCountService().Increase(svcFlowCount)

		appServiceName := constants.AppFlowCountPrefix + app.AppId
		appFlowCount, err := service.GetFlowCountService().GetFlowCount(appServiceName)
		if err != nil {
			return err
		}
		service.GetFlowCountService().Increase(appFlowCount)
		if app.Qpd > 0 && appFlowCount.TotalCount > app.Qpd {
			return err
		}

		if err := handler(srv, stream); err != nil {
			return err
		}

		return nil
	}
}
