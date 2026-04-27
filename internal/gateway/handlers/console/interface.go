package console

import (
	"github.com/vyolayer/vyolayer/pkg/jwt"
	"github.com/vyolayer/vyolayer/pkg/logger"
	consolev1 "github.com/vyolayer/vyolayer/proto/console/v1"
)

type ProjectServiceHandler struct {
	logger *logger.AppLogger
	client consolev1.ProjectServiceManifestClient
	iamJWT jwt.IamJWT
}
