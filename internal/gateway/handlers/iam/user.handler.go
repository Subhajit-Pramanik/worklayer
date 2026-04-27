package iam

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/pkg/errors"
	"github.com/vyolayer/vyolayer/pkg/response"
	pb "github.com/vyolayer/vyolayer/proto/iam/v1"
)

// @Summary Get Current Profile
// @Description Returns the authenticated user's profile
// @Tags IAM Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse{data=GetMeDTO}
// @Router /iam/me [get]
// getMe returns the authenticated user's profile by forwarding to the IAM UserService.
func (h *IAMAuthGatewayHandler) getMe(c *fiber.Ctx) error {

	resp, err := h.user.GetMe(c.UserContext(), &pb.GetMeRequest{})
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	user := resp.GetUser()
	avatar := user.GetAvatar()

	avatarDTO := &AvatarDTO{
		ID:            avatar.GetId(),
		Url:           avatar.GetUrl(),
		FallbackChar:  avatar.GetFallbackChar(),
		FallbackColor: avatar.GetFallbackColor(),
	}

	userDTO := &UserDTO{
		ID:              user.GetId(),
		Email:           user.GetEmail(),
		FullName:        user.GetFullName(),
		Status:          user.GetStatus(),
		IsEmailVerified: user.GetIsEmailVerified(),
		JoinedAt:        user.GetJoinedAt(),
		Avatar:          *avatarDTO,
	}

	respDTO := &GetMeDTO{User: userDTO}
	return response.Success(c, respDTO)
}

// updateMe updates the authenticated user's profile.
func (h *IAMAuthGatewayHandler) updateMe(c *fiber.Ctx) error {

	var req pb.UpdateMeRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, errors.BadRequest("invalid request body"))
	}

	resp, err := h.user.UpdateMe(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusOK, "profile updated", resp.User)
}

// changePassword changes the password for the authenticated user.
func (h *IAMAuthGatewayHandler) changePassword(c *fiber.Ctx) error {

	var req pb.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, errors.BadRequest("invalid request body"))
	}

	if _, err := h.auth.ChangePassword(c.UserContext(), &req); err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusOK, "password changed successfully", nil)
}
