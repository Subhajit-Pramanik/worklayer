package account

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/pkg/errors"
	"github.com/vyolayer/vyolayer/pkg/response"
	accountV1 "github.com/vyolayer/vyolayer/proto/account/v1"
)

// @Summary Register User
// @Description Register a new user account
// @Tags Account
// @Accept json
// @Produce json
// @Param request body RegisterRequestDTO true "Registration details"
// @Success 201 {object} response.SuccessResponse{data=RegisterResponseDTO}
// @Router /account/sign-up [post]
func (h *AccountHandler) register(c *fiber.Ctx) error {
	var req accountV1.RegisterRequest
	if e := c.BodyParser(&req); e != nil {
		return response.Error(c, errors.BadRequest("Invalid Request Body"))
	}

	resp, e := h.client.Register(c.UserContext(), &req)
	if e != nil {
		appErr := errors.FromGRPC(e)
		return response.Error(c, appErr)
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusCreated,
		"User registered successfully",
		resp,
	)
}

// @Summary Verify Email
// @Description Verify a user's email with a token
// @Tags Account
// @Accept json
// @Produce json
// @Param token query string true "Verification token"
// @Success 200 {object} response.SuccessResponse
// @Router /account/verify-email [post]
func (h *AccountHandler) verifyEmail(c *fiber.Ctx) error {
	token := c.Query(QueryParamToken)
	if token == "" {
		return response.Error(c, errors.BadRequest("Token is required"))
	}

	_, e := h.client.VerifyEmail(c.UserContext(), &accountV1.VerifyEmailRequest{
		Token: token,
	})
	if e != nil {
		appErr := errors.FromGRPC(e)
		return response.Error(c, appErr)
	}

	return response.SuccessMessage(
		c,
		"Email verified successfully",
	)
}

// @Summary Resend Verification Email
// @Description Resend the email verification link
// @Tags Account
// @Accept json
// @Produce json
// @Param request body ResendVerificationEmailRequestDTO true "Email details"
// @Success 200 {object} response.SuccessResponse
// @Router /account/resend-verification-email [post]
func (h *AccountHandler) resendVerificationEmail(c *fiber.Ctx) error {
	var req accountV1.ResendVerificationEmailRequest
	if e := c.BodyParser(&req); e != nil {
		return response.Error(c, errors.BadRequest("Invalid Request Body"))
	}

	_, e := h.client.ResendVerificationEmail(c.UserContext(), &req)
	if e != nil {
		return response.Error(c, errors.FromGRPC(e))
	}

	return response.SuccessMessage(
		c,
		"Verification email resent successfully",
	)
}

// @Summary User Login
// @Description Authenticate user and return session tokens
// @Tags Account
// @Accept json
// @Produce json
// @Param request body LoginRequestDTO true "Login credentials"
// @Success 200 {object} response.SuccessResponse{data=LoginResponseDTO}
// @Router /account/sign-in [post]
func (h *AccountHandler) login(c *fiber.Ctx) error {
	var req accountV1.LoginRequest
	if e := c.BodyParser(&req); e != nil {
		return response.Error(c, errors.BadRequest("Invalid Request Body"))
	}

	resp, e := h.client.Login(c.UserContext(), &req)
	if e != nil {
		return response.Error(c, errors.FromGRPC(e))
	}

	if err := h.cookieSv.Set(c, resp.AccessToken, resp.RefreshToken); err != nil {
		return response.Error(c, errors.Internal("Failed to set cookies"))
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"Login successful",
		resp,
	)
}

// @Summary User Logout
// @Description Logout user and clear session cookies
// @Tags Account
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Router /account/sign-out [post]
func (h *AccountHandler) logout(c *fiber.Ctx) error {
	t, err := h.cookieSv.GetRefreshToken(c)
	if err != nil {
		return response.Error(c, errors.BadRequest("Refresh token not found"))
	}

	_, e := h.client.Logout(c.UserContext(), &accountV1.LogoutRequest{
		RefreshToken: t,
	})

	// Always clear the local cookies so the user isn't stuck holding invalid tokens
	if err := h.cookieSv.Clear(c); err != nil {
		log.Printf("Failed to delete refresh token cookie: %v", err)
	}

	if e != nil {
		return response.Error(c, errors.FromGRPC(e))
	}

	return response.SuccessMessage(
		c,
		"Logout successful",
	)
}

// @Summary Validate Session
// @Description Validate the current access token session
// @Tags Account
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Router /account/validate [post]
func (h *AccountHandler) validateSession(c *fiber.Ctx) error {
	// Extract the access token from the cookies
	accessToken := h.cookieSv.GetAccessToken(c)
	if accessToken == "" {
		return response.Error(c, errors.Unauthorized("No access token provided"))
	}

	// Verify the access token securely using the shared JWT verifier
	_, _, err := h.accountJWT.VerifyAccessToken(accessToken)
	if err != nil {
		log.Printf("Session validation failed: %v", err)
		return response.Error(c, errors.Unauthorized("Invalid or expired session"))
	}

	log.Println("Session validated successfully")

	return response.SuccessMessage(
		c,
		"true",
	)
}

// @Summary Refresh Token
// @Description Refresh the access token using the refresh token cookie
// @Tags Account
// @Accept json
// @Produce json
// @Success 200 {object} response.SuccessResponse{data=RefreshSessionResponseDTO}
// @Router /account/sessions/refresh [post]
func (h *AccountHandler) refreshToken(c *fiber.Ctx) error {
	t, err := h.cookieSv.GetRefreshToken(c)
	if err != nil || t == "" {
		return response.Error(c, errors.TokenInvalid("Not found"))
	}

	log.Println("Refreshing session")
	resp, e := h.client.RefreshSession(
		c.UserContext(),
		&accountV1.RefreshSessionRequest{RefreshToken: t},
	)
	if e != nil {
		log.Printf("Failed to refresh session: %v", e)
		return response.Error(c, errors.FromGRPC(e))
	}

	if err := h.cookieSv.Set(
		c,
		resp.AccessToken,
		resp.RefreshToken,
	); err != nil {
		return response.Error(c, errors.Internal("Failed to set cookies"))
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"Session refreshed successfully",
		resp,
	)
}

// @Summary List Sessions
// @Description List all active sessions for the authenticated user
// @Tags Account
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse{data=AllSessionsResponseDTO}
// @Router /account/sessions [get]
func (h *AccountHandler) listSessions(c *fiber.Ctx) error {
	rt, err := h.cookieSv.GetRefreshToken(c)
	if err != nil || rt == "" {
		return response.Error(c, errors.TokenInvalid("Not found"))
	}

	req := &accountV1.AllSessionsRequest{RefreshToken: rt}

	resp, e := h.client.AllSessions(c.UserContext(), req)
	if e != nil {
		return response.Error(c, errors.FromGRPC(e))
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"Sessions listed successfully",
		resp,
	)
}

// @Summary Revoke Session
// @Description Revoke a specific session
// @Tags Account
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param session_id query string true "Session ID"
// @Success 200 {object} response.SuccessResponse
// @Router /account/sessions/revoke [post]
func (h *AccountHandler) revokeSession(c *fiber.Ctx) error {
	// Session id from params
	sessionID := c.Query(QueryParamSessionId)
	if sessionID == "" {
		return response.Error(c, errors.BadRequest("Session id is required"))
	}

	refreshToken, err := h.cookieSv.GetRefreshToken(c)
	if err != nil {
		return response.Error(c, errors.BadRequest("Refresh token not found"))
	}

	_, e := h.client.RevokeSession(c.UserContext(), &accountV1.RevokeSessionRequest{
		RefreshToken: refreshToken,
		SessionId:    sessionID,
	})
	if e != nil {
		return response.Error(c, errors.FromGRPC(e))
	}

	if err := h.cookieSv.Clear(c); err != nil {
		log.Printf("Failed to delete refresh token cookie: %v", err)
	}

	return response.SuccessMessage(
		c,
		"Session revoked successfully",
	)
}

// @Summary Revoke All Sessions
// @Description Revoke all sessions for the authenticated user
// @Tags Account
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Router /account/sessions/revoke-all [post]
func (h *AccountHandler) revokeAllSessions(c *fiber.Ctx) error {
	refreshToken, err := h.cookieSv.GetRefreshToken(c)
	if err != nil || refreshToken == "" {
		return response.Error(c, errors.BadRequest("Refresh token not found"))
	}

	_, e := h.client.RevokeAllSessions(c.UserContext(), &accountV1.RevokeAllSessionsRequest{
		RefreshToken: refreshToken,
	})
	if e != nil {
		return response.Error(c, errors.FromGRPC(e))
	}

	if err := h.cookieSv.Clear(c); err != nil {
		log.Printf("Failed to delete refresh token cookie: %v", err)
	}

	return response.SuccessMessage(
		c,
		"All sessions revoked successfully",
	)
}

// Account recover - Change Password, Forgot Password, Reset Password

// @Summary Change Password
// @Description Change the user's password
// @Tags Account
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChangePasswordRequestDTO true "Change Password details"
// @Success 200 {object} response.SuccessResponse
// @Router /account/change-password [post]
func (h *AccountHandler) changePassword(c *fiber.Ctx) error {
	var req accountV1.ChangePasswordRequest
	if e := c.BodyParser(&req); e != nil {
		return response.Error(c, errors.BadRequest("Invalid Request Body"))
	}

	_, err := h.client.ChangePassword(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	// clear cookies
	h.cookieSv.Clear(c)

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"Password changed successfully",
		nil,
	)
}

// @Summary Forgot Password
// @Description Send a password reset email
// @Tags Account
// @Accept json
// @Produce json
// @Param request body ForgotPasswordRequestDTO true "Forgot Password details"
// @Success 200 {object} response.SuccessResponse
// @Router /account/forgot-password [post]
func (h *AccountHandler) forgotPassword(c *fiber.Ctx) error {
	var req accountV1.ForgotPasswordRequest
	if e := c.BodyParser(&req); e != nil {
		return response.Error(c, errors.BadRequest("Invalid Request Body"))
	}

	_, err := h.client.ForgotPassword(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"Email sent successfully",
		nil,
	)
}

// @Summary Reset Password
// @Description Reset the user's password using a token
// @Tags Account
// @Accept json
// @Produce json
// @Param token query string true "Reset token"
// @Param request body ResetPasswordRequestDTO true "Reset Password details"
// @Success 200 {object} response.SuccessResponse
// @Router /account/reset-password [post]
func (h *AccountHandler) resetPassword(c *fiber.Ctx) error {
	token := c.Query(QueryParamToken)
	if token == "" {
		return response.Error(c, errors.BadRequest("Token is required"))
	}

	var req accountV1.ResetPasswordRequest
	if e := c.BodyParser(&req); e != nil {
		return response.Error(c, errors.BadRequest("Invalid Request Body"))
	}
	req.Token = token

	_, err := h.client.ResetPassword(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	h.cookieSv.Clear(c)

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"Password reset successfully",
		nil,
	)
}
