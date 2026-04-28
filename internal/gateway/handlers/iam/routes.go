package iam

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/internal/gateway/middleware"
	"github.com/vyolayer/vyolayer/internal/gateway/service"
	"github.com/vyolayer/vyolayer/pkg/logger"
	iAMV1 "github.com/vyolayer/vyolayer/proto/iam/v1"
)

func NewIAMAuthGatewayHandler(
	auth iAMV1.AuthServiceClient,
	user iAMV1.UserServiceClient,
	cookie *service.IAMCookieService,
	authClient iAMV1.AuthServiceClient,
	logger *logger.AppLogger,
) *IAMAuthGatewayHandler {
	return &IAMAuthGatewayHandler{
		auth:       auth,
		user:       user,
		cookie:     cookie,
		authClient: authClient,
		logger:     logger.WithContext("IAMAuthGatewayHandler"),
	}
}

// RegisterRoutes mounts all IAM routes under /iam.
func (h *IAMAuthGatewayHandler) RegisterRoutes(router fiber.Router) {
	// grpc ctx timeout
	grpcCtxMiddleware := middleware.NewGrpcCtxMiddleware(grpcTimeout)

	registerLimiter := middleware.NewRateLimiter(10, 1*time.Minute, "email")
	emailVerifyLimiter := middleware.NewRateLimiter(10, 1*time.Minute, "email-verify")
	loginLimiter := middleware.NewRateLimiter(10, 1*time.Minute, "login")

	standardLimiter := middleware.NewRateLimiter(100, 1*time.Minute, "auth").Handler()

	iam := router.Group("/iam")
	iam.Use(grpcCtxMiddleware.Handler())

	// ── Public auth endpoints ────────────────────────────────────────────────
	iam.Post("/register", registerLimiter.Handler(), h.register)
	iam.Post("/verify-email", standardLimiter, h.verifyEmail)
	iam.Post("/resend-verification-email", emailVerifyLimiter.Handler(), h.resendVerificationEmail)

	iam.Post("/login", loginLimiter.Handler(), h.login)
	iam.Post("/logout", standardLimiter, h.logout)
	iam.Post("/refresh-session", standardLimiter, h.refreshSession)

	// iam.Post("/forgot-password", strictLimiter, h.forgotPassword)
	// iam.Post("/reset-password", strictLimiter, h.resetPassword)

	// ── Authenticated profile endpoints (/me) ───────────────────────────────
	me := iam.Group("/me", standardLimiter)
	me.Use(middleware.IamJWTVerify(h.authClient))
	me.Get("/", h.getMe)
	// me.Patch("/", h.updateMe)
	// me.Post("/change-password", h.changePassword)

	h.logger.Info("IAM routes registered", "")
}
