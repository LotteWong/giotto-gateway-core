package grpc_proxy_middleware

import (
	"fmt"
	"strings"

	"github.com/LotteWong/giotto-gateway/constants"
	"github.com/LotteWong/giotto-gateway/models/po"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func GrpcIpAuthMiddleware(grpcServiceDetail *po.ServiceDetail) func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		var whiteIpList []string
		var blackIpList []string
		openAuth := grpcServiceDetail.AccessControl.OpenAuth
		whiteList := grpcServiceDetail.AccessControl.WhiteList
		if whiteList != "" {
			whiteIpList = strings.Split(whiteList, ",")
		}
		blackList := grpcServiceDetail.AccessControl.BlackList
		if blackList != "" {
			blackIpList = strings.Split(blackList, ",")
		}

		if openAuth == constants.Enable {
			peerCtx, ok := peer.FromContext(stream.Context())
			if !ok {
				return errors.New("failed to get peer context")
			}
			peerAddr := peerCtx.Addr.String()
			clientIp := peerAddr[0:strings.LastIndex(peerAddr, ":")]

			if len(whiteIpList) > 0 { // white list has higher priority
				if !checkStrInSlice(whiteIpList, clientIp) {
					return errors.New(fmt.Sprintf("ip %s not in white ip list", clientIp))
				}
			} else { // black list has lower priority
				if len(blackIpList) > 0 {
					if checkStrInSlice(blackIpList, clientIp) {
						return errors.New(fmt.Sprintf("ip %s is in black ip list", clientIp))
					}
				}
			}
		}

		if err := handler(srv, stream); err != nil {
			return err
		}

		return nil
	}
}

func checkStrInSlice(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}
