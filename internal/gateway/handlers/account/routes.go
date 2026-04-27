package account

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/internal/gateway/middleware"
	"github.com/vyolayer/vyolayer/internal/gateway/service"
	sharedMiddleware "github.com/vyolayer/vyolayer/internal/shared/middleware"
	"github.com/vyolayer/vyolayer/pkg/jwt"
	"github.com/vyolayer/vyolayer/pkg/logger"
	accountV1 "github.com/vyolayer/vyolayer/proto/account/v1"
)

// NewAccountHandler creates a new AccountHandler injecting the gRPC client
func NewAccountHandler(
	client accountV1.AccountServiceClient,
	cookieSv *service.AccountTokenService,
	accountJWT jwt.AccountJWT,
	logger *logger.AppLogger,
) *AccountHandler {
	return &AccountHandler{
		client:     client,
		cookieSv:   cookieSv,
		accountJWT: accountJWT,
		logger:     logger.WithContext("AccountHandler"),
	}
}

// RegisterRoutes registers the account routes on the provided router
func (h *AccountHandler) RegisterRoutes(router fiber.Router) {
	grpcCtxMiddleware := middleware.NewGrpcCtxMiddleware(grpcTimeout)

	r := router.Group("/account")
	r.Use(grpcCtxMiddleware.Handler())

	r.Post("/sign-up", h.register)
	r.Post("/verify-email", h.verifyEmail)
	r.Post("/resend-verification-email", h.resendVerificationEmail)
	r.Post("/sign-in", h.login)

	r.Post("/sessions/refresh", h.refreshToken)

	r.Post("/forgot-password", h.forgotPassword)
	r.Post("/reset-password", h.resetPassword)

	ra := r.Group("/")
	ra.Use(sharedMiddleware.AccountJWTVerify(h.accountJWT))

	ra.Post("/sign-out", h.logout)
	ra.Post("/validate", h.validateSession)

	ra.Get("/sessions", h.listSessions)
	ra.Post("/sessions/revoke", h.revokeSession)
	ra.Post("/sessions/revoke-all", h.revokeAllSessions)

	ra.Post("/change-password", h.changePassword)

	h.logger.Info("Account routes registered", "")
}
