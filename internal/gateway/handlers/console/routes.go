package console

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/internal/gateway/middleware"
	"github.com/vyolayer/vyolayer/pkg/jwt"
	"github.com/vyolayer/vyolayer/pkg/logger"
	consolev1 "github.com/vyolayer/vyolayer/proto/console/v1"
)

func NewProjectServiceHandler(
	logger *logger.AppLogger,
	client consolev1.ProjectServiceManifestClient,
	iamJWT jwt.IamJWT,
) *ProjectServiceHandler {
	return &ProjectServiceHandler{
		logger: logger.WithContext("ConsoleProjectServiceHandler"),
		client: client,
		iamJWT: iamJWT,
	}
}

func (h *ProjectServiceHandler) RegisterRoutes(router fiber.Router) {
	grpcCtxMiddleware := middleware.NewGrpcCtxMiddleware(grpcTimeout)

	services := router.Group("/console/projects/:" + ParamProjectID + "/services")
	services.Use(grpcCtxMiddleware.Handler())
	services.Use(middleware.IamJWTVerify(h.iamJWT))
	services.Use(middleware.ValidateProjectID())

	services.Get("/", h.list)
	services.Get("/:"+ParamServiceKey, h.get)

	h.logger.Info("ProjectService routes registered", "")
}
