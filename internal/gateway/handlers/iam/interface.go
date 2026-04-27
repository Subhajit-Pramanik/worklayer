package iam

import (
	"github.com/vyolayer/vyolayer/internal/gateway/service"
	"github.com/vyolayer/vyolayer/pkg/jwt"
	"github.com/vyolayer/vyolayer/pkg/logger"
	iAMV1 "github.com/vyolayer/vyolayer/proto/iam/v1"
)

type IAMAuthGatewayHandler struct {
	auth   iAMV1.AuthServiceClient
	user   iAMV1.UserServiceClient
	cookie *service.IAMCookieService
	iamJWT jwt.IamJWT
	logger *logger.AppLogger
}
