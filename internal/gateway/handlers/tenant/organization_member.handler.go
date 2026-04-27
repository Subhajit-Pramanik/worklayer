package tenant

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/pkg/errors"
	"github.com/vyolayer/vyolayer/pkg/response"
	tenantV1 "github.com/vyolayer/vyolayer/proto/tenant/v1"
)

// @Summary Get Current Organization Member
// @Description Get details of the currently authenticated member
// @Tags Organization Members
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param organizationID path string true "Organization ID"
// @Success 200 {object} response.SuccessResponse{data=OrganizationMemberWithRBACResponse}
// @Router /organizations/{organizationID}/members/me [get]
func (h *OrganizationMemberHandler) getCurrentMember(c *fiber.Ctx) error {
	resp, err := h.client.GetCurrentMember(
		c.UserContext(),
		&tenantV1.TenantOrganizationIDRequest{OrganizationId: getOrgIDFromLocals(c)},
	)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	member := resp.GetMember()
	if member == nil {
		return response.Error(c, errors.NotFound("member not found"))
	}

	permsDto := make([]string, len(resp.GetPermissions()))
	for i, p := range resp.GetPermissions() {
		permsDto[i] = p.GetCode()
	}

	rolesDto := make([]string, len(resp.GetRoles()))
	for i, r := range resp.GetRoles() {
		rolesDto[i] = r.GetName()
	}

	memberDto := &OrganizationMemberWithRBACResponse{
		OrganizationMember: *protoMemberToDTO(member),
		Roles:              rolesDto,
		Perms:              permsDto,
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"member fetched successfully",
		memberDto,
	)
}

// @Summary List Organization Members
// @Description List all members in the organization
// @Tags Organization Members
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param organizationID path string true "Organization ID"
// @Success 200 {object} response.SuccessResponse{data=ListOrganizationMembersResponse}
// @Router /organizations/{organizationID}/members [get]
func (h *OrganizationMemberHandler) listMembers(c *fiber.Ctx) error {
	req := tenantV1.TenantOrganizationIDRequest{
		OrganizationId: getOrgIDFromLocals(c),
	}

	resp, err := h.client.GetAllMembersByOrg(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	membersDto := make([]*OrganizationMember, len(resp.GetMembers()))
	for i, m := range resp.GetMembers() {
		membersDto[i] = protoMemberToDTO(m)
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"members fetched successfully",
		&ListOrganizationMembersResponse{
			Members:    membersDto,
			TotalCount: resp.GetTotalCount(),
		},
	)
}

// @Summary Get Organization Member
// @Description Get details of a specific organization member
// @Tags Organization Members
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param organizationID path string true "Organization ID"
// @Param memberID path string true "Member ID"
// @Success 200 {object} response.SuccessResponse{data=OrganizationMember}
// @Router /organizations/{organizationID}/members/{memberID} [get]
func (h *OrganizationMemberHandler) getMemberByID(c *fiber.Ctx) error {
	req := tenantV1.GetOrganizationMemberByIdRequest{
		OrganizationId: getOrgIDFromLocals(c),
		MemberId:       c.Params(ParamMemberID),
	}

	resp, err := h.client.GetMemberById(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	member := resp.GetMember()
	if member == nil {
		return response.Error(c, errors.NotFound("member not found"))
	}

	memberDto := protoMemberToDTO(member)

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"member fetched successfully",
		memberDto,
	)
}

// @Summary Remove Organization Member
// @Description Remove a member from the organization
// @Tags Organization Members
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param organizationID path string true "Organization ID"
// @Param memberID path string true "Member ID"
// @Success 200 {object} response.SuccessResponse
// @Router /organizations/{organizationID}/members/{memberID} [delete]
func (h *OrganizationMemberHandler) removeMember(c *fiber.Ctx) error {
	req := tenantV1.RemoveOrganizationMemberRequest{
		OrganizationId: getOrgIDFromLocals(c),
		MemberId:       c.Params(ParamMemberID),
	}

	resp, err := h.client.RemoveMember(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusOK, resp.GetMessage(), nil)
}
