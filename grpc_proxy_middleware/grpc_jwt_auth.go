package grpc_proxy_middleware

import (
	"fmt"
	"strings"

	"github.com/LotteWong/giotto-gateway-core/constants"
	"github.com/LotteWong/giotto-gateway-core/models/po"
	"github.com/LotteWong/giotto-gateway-core/service"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func GrpcJwtAuthMiddleware(grpcServiceDetail *po.ServiceDetail) func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		var jwtErr error
		var hdlErr error

		metaCtx, ok := metadata.FromIncomingContext(stream.Context())
		if !ok {
			return errors.New("failed to get metadata context")
		}

		// parse authorization to get jwt
		tokens := metaCtx.Get("Authorization")
		if len(tokens) == 0 {
			return errors.New("auth not found")
		}
		pair := strings.Split(tokens[0], " ")
		if len(pair) != 2 {
			return errors.New("can not get jwt from authorization header")
		}

		tokenType := pair[0]
		tokenString := pair[1]

		// verify jwt by expire at and issuer
		switch tokenType {
		case constants.JwtType:
			jwtErr = service.GetJwtService().GrpcVerifyJwt(metaCtx, grpcServiceDetail, tokenString)
		default:
			jwtErr = errors.New(fmt.Sprintf("not support jwt type %s", tokenType))
		}
		if jwtErr != nil {
			return jwtErr
		}

		hdlErr = handler(srv, stream)
		if hdlErr != nil {
			return hdlErr
		}

		return nil
	}
}
