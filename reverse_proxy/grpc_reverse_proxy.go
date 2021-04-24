package reverse_proxy

import (
	"context"
	"log"

	"github.com/LotteWong/giotto-gateway-core/load_balance"
	"github.com/e421083458/grpc-proxy/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func NewGrpcReverseProxy(lb load_balance.LoadBalance) grpc.StreamHandler {
	return func() grpc.StreamHandler {
		addr, err := lb.Get("")
		if err != nil {
			log.Fatal("failed to get addr")
		}
		director := func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
			c, err := grpc.DialContext(ctx, addr, grpc.WithCodec(proxy.Codec()), grpc.WithInsecure())
			md, _ := metadata.FromIncomingContext(ctx)
			outCtx, _ := context.WithCancel(ctx)
			outCtx = metadata.NewOutgoingContext(outCtx, md.Copy())
			return outCtx, c, err
		}
		return proxy.TransparentHandler(director)
	}()
}
