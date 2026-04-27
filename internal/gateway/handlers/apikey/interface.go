package apikey

import (
	"github.com/vyolayer/vyolayer/pkg/jwt"
	"github.com/vyolayer/vyolayer/pkg/logger"
	pb "github.com/vyolayer/vyolayer/proto/apikey/v1"
)

type ApiKeyHandler struct {
	logger *logger.AppLogger
	client pb.APIKeyServiceClient
	iamJWT jwt.IamJWT
}
