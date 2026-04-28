package console

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/internal/gateway/middleware"
	"github.com/vyolayer/vyolayer/pkg/logger"
	consolev1 "github.com/vyolayer/vyolayer/proto/console/v1"
	iAMV1 "github.com/vyolayer/vyolayer/proto/iam/v1"
)

func NewProjectServiceHandler(
	logger *logger.AppLogger,
	client consolev1.ProjectServiceManifestClient,
	authClient iAMV1.AuthServiceClient,

) *ProjectServiceHandler {
	return &ProjectServiceHandler{
		logger:     logger.WithContext("ConsoleProjectServiceHandler"),
		client:     client,
		authClient: authClient,
	}
}

func (h *ProjectServiceHandler) RegisterRoutes(router fiber.Router) {
	grpcCtxMiddleware := middleware.NewGrpcCtxMiddleware(grpcTimeout)

	services := router.Group("/console/projects/:" + ParamProjectID + "/services")
	services.Use(grpcCtxMiddleware.Handler())
	services.Use(middleware.IamJWTVerify(h.authClient))
	services.Use(middleware.ValidateProjectID())

	services.Get("/", h.list)
	services.Get("/:"+ParamServiceKey, h.get)

	h.logger.Info("ProjectService routes registered", "")
}
