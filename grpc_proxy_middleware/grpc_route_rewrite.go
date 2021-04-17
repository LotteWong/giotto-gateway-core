package grpc_proxy_middleware

import (
	"log"
	"strings"

	"github.com/LotteWong/giotto-gateway/models/po"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func GrpcRouteRewriteMiddleware(grpcServiceDetail *po.ServiceDetail) func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		metaCtx, ok := metadata.FromIncomingContext(stream.Context())
		if !ok {
			return errors.New("failed to get metadata context")
		}

		for _, rule := range strings.Split(grpcServiceDetail.GrpcRule.HeaderTransform, ",") {
			items := strings.Split(rule, " ")
			if len(items) < 2 || len(items) > 3 {
				log.Println("url rewrite format error")
				continue
			}
			action := items[0]

			if len(items) == 2 {
				key := items[1]
				if action == "delete" {
					delete(metaCtx, key)
				}
			} else {
				key, value := items[1], items[2]
				if action == "add" || action == "edit" {
					metaCtx.Set(key, value)
				}
			}

		}

		if err := stream.SetHeader(metaCtx); err != nil {
			return err
		}

		if err := handler(srv, stream); err != nil {
			return err
		}

		return nil
	}
}
