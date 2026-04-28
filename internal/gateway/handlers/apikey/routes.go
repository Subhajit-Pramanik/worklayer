package apikey

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/internal/gateway/middleware"
	"github.com/vyolayer/vyolayer/pkg/logger"
	pb "github.com/vyolayer/vyolayer/proto/apikey/v1"
	iAMV1 "github.com/vyolayer/vyolayer/proto/iam/v1"
)

func NewHandler(
	logger *logger.AppLogger,
	client pb.APIKeyServiceClient,
	authClient iAMV1.AuthServiceClient,

) *ApiKeyHandler {
	return &ApiKeyHandler{
		logger:     logger.WithContext("APIKeyHandler"),
		client:     client,
		authClient: authClient,
	}
}

func (h *ApiKeyHandler) RegisterRoutes(router fiber.Router) {
	grpcCtxMiddleware := middleware.NewGrpcCtxMiddleware(grpcTimeout).Handler()

	apiKey := router.Group("/api-keys")
	apiKey.Use(grpcCtxMiddleware, middleware.IamJWTVerify(h.authClient))

	apiKey.Get("/", h.List)
	apiKey.Post("/", h.Create)

	apiKey.Get("/:apiKeyID", h.Get)
	apiKey.Patch("/:apiKeyID/rotate", h.Rotate)
	apiKey.Delete("/:apiKeyID/revoke", h.Revoke)
	apiKey.Get("/:apiKeyID/validate", h.Validate)

	h.logger.Info("APIKey routes registered", "")
}
