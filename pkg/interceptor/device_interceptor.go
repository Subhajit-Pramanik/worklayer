package interceptor

import (
	"context"

	"github.com/vyolayer/vyolayer/pkg/ctxutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func DeviceInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		ip := getIP(md)
		userAgent := getUserAgent(md)

		ctx = ctxutil.InjectDeviceInfo(ctx, ip, userAgent)
		return handler(ctx, req)
	}
}

func getIP(md metadata.MD) string {
	ip := ""
	if len(md.Get("x-forwarded-for")) > 0 {
		ip = md.Get("x-forwarded-for")[0]
	} else if len(md.Get("x-real-ip")) > 0 {
		ip = md.Get("x-real-ip")[0]
	}

	return ip
}

func getUserAgent(md metadata.MD) string {
	userAgent := ""
	if len(md.Get("user-agent")) > 0 {
		userAgent = md.Get("user-agent")[0]
	}

	return userAgent
}
