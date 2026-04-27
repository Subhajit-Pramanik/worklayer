package iam

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/pkg/errors"
	"github.com/vyolayer/vyolayer/pkg/response"

	pb "github.com/vyolayer/vyolayer/proto/iam/v1"
)

// @Summary Register IAM User
// @Description Register a new user in IAM
// @Tags IAM Auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} response.SuccessResponse
// @Router /iam/register [post]
func (h *IAMAuthGatewayHandler) register(c *fiber.Ctx) error {

	var (
		ctx = c.UserContext()
		in  RegisterRequest
	)

	if err := c.BodyParser(&in); err != nil {
		return response.Error(c, errors.BadRequest("invalid request body"))
	}

	if _, err := h.auth.Register(ctx, &pb.RegisterRequest{
		FullName: in.FullName,
		Email:    in.Email,
		Password: in.Password,
	}); err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusCreated, "User registered successfully", nil)
}

// @Summary Verify Email
// @Description Verify a user's email with a token
// @Tags IAM Auth
// @Accept json
// @Produce json
// @Param token query string true "Verification token"
// @Success 200 {object} response.SuccessResponse
// @Router /iam/verify-email [post]
func (h *IAMAuthGatewayHandler) verifyEmail(c *fiber.Ctx) error {

	var (
		ctx   = c.UserContext()
		token string
	)

	if token = c.Query(QueryParamToken); token == "" {
		return response.Error(c, errors.BadRequest("token is required"))
	}

	if _, err := h.auth.VerifyEmail(ctx, &pb.VerifyEmailRequest{
		Token: token,
	}); err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusOK, "Email verified successfully", nil)
}

// @Summary Resend Verification Email
// @Description Resend the email verification link
// @Tags IAM Auth
// @Accept json
// @Produce json
// @Param request body ResendVerificationEmailRequest true "Email details"
// @Success 200 {object} response.SuccessResponse
// @Router /iam/resend-verification-email [post]
func (h *IAMAuthGatewayHandler) resendVerificationEmail(c *fiber.Ctx) error {

	var (
		ctx = c.UserContext()
		req ResendVerificationEmailRequest
	)

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, errors.BadRequest("invalid request body"))
	}

	if _, err := h.auth.ResendVerificationEmail(ctx, &pb.ResendVerificationEmailRequest{
		Email: req.Email,
	}); err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusOK, "Verification email resent successfully", nil)
}

// ── Session ─────────────────────────────────────────────────────────────────────

// @Summary User Login
// @Description Authenticate user and return session tokens
// @Tags IAM Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} response.SuccessResponse{data=AuthSessionDTO}
// @Router /iam/login [post]
func (h *IAMAuthGatewayHandler) login(c *fiber.Ctx) error {

	var (
		ctx = c.UserContext()
		req LoginRequest
	)

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, errors.BadRequest("invalid request body"))
	}

	sess, err := h.auth.Login(ctx, &pb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	if err := h.cookie.Set(
		c,
		sess.AccessToken,
		sess.SessionToken,
		sess.AccessTokenExpiresAt.AsTime(),
		sess.SessionTokenExpiresAt.AsTime(),
	); err != nil {
		return response.Error(c, errors.Internal("failed to set cookies"))
	}

	resp := AuthSessionDTO{
		AccessToken:          sess.AccessToken,
		AccessTokenExpiresAt: sess.AccessTokenExpiresAt.AsTime().Unix(),
	}

	return response.Send(c, fiber.StatusOK, "Login successful", resp)
}

// @Summary User Logout
// @Description Logout user and clear session cookies
// @Tags IAM Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Router /iam/logout [post]
func (h *IAMAuthGatewayHandler) logout(c *fiber.Ctx) error {

	var (
		ctx = c.UserContext()
		st  string
	)

	if st = h.cookie.GetSessionToken(c); st == "" {
		return response.Error(c, errors.Unauthorized("unauthorized"))
	}

	if _, err := h.auth.Logout(ctx, &pb.LogoutRequest{SessionToken: st}); err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	if err := h.cookie.Clear(c); err != nil {
		log.Printf("[IAM] failed to clear cookies: %v", err)
		return response.Error(c, errors.Internal("failed to clear cookies"))
	}

	return response.SuccessWithMessage(c, fiber.StatusOK, "logged out successfully", nil)
}

// @Summary Refresh Session
// @Description Refresh the access token session
// @Tags IAM Auth
// @Accept json
// @Produce json
// @Success 200 {object} response.SuccessResponse
// @Router /iam/refresh-session [post]
func (h *IAMAuthGatewayHandler) refreshSession(c *fiber.Ctx) error {

	var (
		ctx = c.UserContext()
		st  string
	)

	if st = h.cookie.GetSessionToken(c); st == "" {
		return response.Error(c, errors.Unauthorized("unauthorized"))
	}

	sess, err := h.auth.RefreshSession(ctx, &pb.RefreshSessionRequest{SessionToken: st})
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	if err := h.cookie.Set(
		c,
		sess.AccessToken,
		sess.SessionToken,
		sess.AccessTokenExpiresAt.AsTime(),
		sess.SessionTokenExpiresAt.AsTime(),
	); err != nil {
		return response.Error(c, errors.Internal("failed to set cookies"))
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"session refreshed",
		fiber.Map{"access_token": sess.AccessToken},
	)
}

// ── Password ─────────────────────────────────────────────────────────────────────

func (h *IAMAuthGatewayHandler) forgotPassword(c *fiber.Ctx) error {

	var req pb.ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, errors.BadRequest("invalid request body"))
	}

	if _, err := h.auth.ForgotPassword(c.UserContext(), &req); err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	// Always return 200 to avoid leaking whether the email exists.
	return response.SuccessWithMessage(c, fiber.StatusOK, "if this email is registered, a reset link has been sent", nil)
}

func (h *IAMAuthGatewayHandler) resetPassword(c *fiber.Ctx) error {

	var req pb.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, errors.BadRequest("invalid request body"))
	}

	if _, err := h.auth.ResetPassword(c.UserContext(), &req); err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusOK, "password reset successfully", nil)
}
