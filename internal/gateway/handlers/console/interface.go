package console

import (
	"github.com/vyolayer/vyolayer/pkg/logger"
	consolev1 "github.com/vyolayer/vyolayer/proto/console/v1"
	iAMV1 "github.com/vyolayer/vyolayer/proto/iam/v1"
)

type ProjectServiceHandler struct {
	logger     *logger.AppLogger
	client     consolev1.ProjectServiceManifestClient
	authClient iAMV1.AuthServiceClient
}
