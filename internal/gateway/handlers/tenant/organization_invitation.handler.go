package tenant

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/pkg/ctxutil"
	"github.com/vyolayer/vyolayer/pkg/errors"
	"github.com/vyolayer/vyolayer/pkg/response"
	tenantV1 "github.com/vyolayer/vyolayer/proto/tenant/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// @Summary Create Organization Invitation
// @Description Invite a user to an organization
// @Tags Organization Invitations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param organizationID path string true "Organization ID"
// @Param request body tenantV1.CreateInvitationRequest true "Invitation details"
// @Success 201 {object} response.SuccessResponse
// @Router /organizations/{organizationID}/invitations [post]
func (h *OrganizationInvitationHandler) createInvitation(c *fiber.Ctx) error {

	var req tenantV1.CreateInvitationRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, ErrInvalidBody)
	}
	req.OrganizationId = getOrgIDFromLocals(c)

	resp, err := h.client.CreateInvitation(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusCreated, resp.GetMessage(), nil)
}

// @Summary List Organization Invitations
// @Description List invitations for an organization
// @Tags Organization Invitations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param organizationID path string true "Organization ID"
// @Param page_size query int false "Page size"
// @Param page_token query string false "Page token"
// @Success 200 {object} response.SuccessResponse{data=ListOrganizationInvitationsResponse}
// @Router /organizations/{organizationID}/invitations [get]
func (h *OrganizationInvitationHandler) listInvitations(c *fiber.Ctx) error {

	req := tenantV1.ListInvitationsRequest{
		OrganizationId: getOrgIDFromLocals(c),
		PageSize:       int32(c.QueryInt(QueryParamPageSize, 0)),
		PageToken:      c.Query(QueryParamPageToken, ""),
	}

	resp, err := h.client.ListInvitations(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	invitationsDto := make([]*OrganizationInvitation, len(resp.GetInvitations()))
	for i, inv := range resp.GetInvitations() {
		invitationsDto[i] = protoInvitationToDTO(inv)
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"invitations fetched successfully",
		&ListOrganizationInvitationsResponse{Invitations: invitationsDto},
	)
}

// @Summary Cancel Organization Invitation
// @Description Cancel a pending invitation by its ID
// @Tags Organization Invitations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param organizationID path string true "Organization ID"
// @Param invitationID path string true "Invitation ID"
// @Success 200 {object} response.SuccessResponse
// @Router /organizations/{organizationID}/invitations/{invitationID} [delete]
// cancelInvitation cancels a pending invitation by its ID.
func (h *OrganizationInvitationHandler) cancelInvitation(c *fiber.Ctx) error {
	req := tenantV1.CancelInvitationRequest{
		OrganizationId: getOrgIDFromLocals(c),
		InvitationId:   c.Params(ParamInvitationID),
	}

	resp, err := h.client.CancelInvitation(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusOK, resp.GetMessage(), nil)
}

// @Summary Get Pending Invitations by Organization
// @Description Get all pending invitations for the specified organization
// @Tags Organization Invitations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param organizationID path string true "Organization ID"
// @Success 200 {object} response.SuccessResponse{data=ListOrganizationInvitationsForOrgResponse}
// @Router /organizations/{organizationID}/invitations/pending [get]
// getPendingInvitationsByOrgID returns all pending invitations for the specified orgID.
func (h *OrganizationInvitationHandler) getPendingByOrgID(c *fiber.Ctx) error {
	req := tenantV1.TenantOrganizationIDRequest{
		OrganizationId: getOrgIDFromLocals(c),
	}

	resp, err := h.client.GetPendingByOrg(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	invitationsDto := make([]*OrganizationInvitationForOrg, len(resp.GetInvitations()))
	for i, inv := range resp.GetInvitations() {
		invitationsDto[i] = protoInvitationForOrgToDTO(inv)
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"pending invitations fetched successfully",
		&ListOrganizationInvitationsForOrgResponse{Invitations: invitationsDto},
	)
}

// @Summary Get Pending Invitations for User
// @Description Get all pending invitations for the authenticated user
// @Tags Organization Invitations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse{data=ListOrganizationInvitationsResponse}
// @Router /organizations/invitations/pending [get]
// getPendingInvitations returns all pending invitations for the authenticated user.
// This bypasses the org context check – it is scoped to the calling user.
func (h *OrganizationInvitationHandler) getPendingByUser(c *fiber.Ctx) error {
	userEmail, err := ctxutil.ExtractIAMUserEmail(c.UserContext())
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	req := tenantV1.GetPendingInvitationsRequest{
		Email: userEmail,
	}

	resp, err := h.client.GetPendingInvitations(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	invitationsDto := make([]*OrganizationInvitation, len(resp.GetInvitations()))
	for i, inv := range resp.GetInvitations() {
		invitationsDto[i] = protoInvitationToDTO(inv)
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"pending invitations fetched successfully",
		&ListOrganizationInvitationsResponse{Invitations: invitationsDto},
	)
}

// @Summary Accept Organization Invitation
// @Description Accept an invitation using the token
// @Tags Organization Invitations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param token query string true "Invitation token"
// @Success 200 {object} response.SuccessResponse
// @Router /organizations/invitations/accept [post]
// acceptInvitation accepts an invitation using the token from the request body.
// It is unauthenticated at the org level (user may not be a member yet).
func (h *OrganizationInvitationHandler) acceptInvitation(c *fiber.Ctx) error {
	var req tenantV1.AcceptInvitationRequest
	token := c.Query(QueryParamToken)
	if token == "" {
		return response.Error(c, status.Error(codes.InvalidArgument, "token is required"))
	}
	req.Token = token

	resp, err := h.client.AcceptInvitation(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusOK, resp.GetMessage(), nil)
}
