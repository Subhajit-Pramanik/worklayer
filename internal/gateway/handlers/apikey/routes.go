package apikey

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/internal/gateway/middleware"
	"github.com/vyolayer/vyolayer/pkg/jwt"
	"github.com/vyolayer/vyolayer/pkg/logger"
	pb "github.com/vyolayer/vyolayer/proto/apikey/v1"
)

func NewHandler(
	logger *logger.AppLogger,
	client pb.APIKeyServiceClient,
	iamJWT jwt.IamJWT,
) *ApiKeyHandler {
	return &ApiKeyHandler{
		logger: logger.WithContext("APIKeyHandler"),
		client: client,
		iamJWT: iamJWT,
	}
}

func (h *ApiKeyHandler) RegisterRoutes(router fiber.Router) {
	grpcCtxMiddleware := middleware.NewGrpcCtxMiddleware(grpcTimeout)

	apiKey := router.Group("/api-keys")
	apiKey.Use(grpcCtxMiddleware.Handler())
	apiKey.Use(middleware.IamJWTVerify(h.iamJWT))

	apiKey.Get("/", h.List)
	apiKey.Post("/", h.Create)

	apiKey.Get("/:apiKeyID", h.Get)
	apiKey.Patch("/:apiKeyID/rotate", h.Rotate)
	apiKey.Delete("/:apiKeyID/revoke", h.Revoke)
	apiKey.Get("/:apiKeyID/validate", h.Validate)

	h.logger.Info("APIKey routes registered", "")
}
