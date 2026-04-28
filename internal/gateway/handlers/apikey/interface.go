package apikey

import (
	"github.com/vyolayer/vyolayer/pkg/logger"
	pb "github.com/vyolayer/vyolayer/proto/apikey/v1"
	iAMV1 "github.com/vyolayer/vyolayer/proto/iam/v1"
)

type ApiKeyHandler struct {
	logger     *logger.AppLogger
	client     pb.APIKeyServiceClient
	authClient iAMV1.AuthServiceClient
}
