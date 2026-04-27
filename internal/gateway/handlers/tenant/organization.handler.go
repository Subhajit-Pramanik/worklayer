package tenant

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/pkg/errors"
	"github.com/vyolayer/vyolayer/pkg/response"
	tenantV1 "github.com/vyolayer/vyolayer/proto/tenant/v1"
)

// ─── Organization ────────────────────────────────────────────────────────────

// @Summary Create Organization
// @Description Create a new organization
// @Tags Organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body tenantV1.CreateOrganizationRequest true "Organization details"
// @Success 201 {object} response.SuccessResponse{data=CreateOrganizationResponse}
// @Router /organizations [post]
func (h *OrganizationHandler) create(c *fiber.Ctx) error {
	var req tenantV1.CreateOrganizationRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, ErrInvalidBody)
	}

	resp, err := h.client.CreateOrganization(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	h.logger.Debug("Organization created", resp)

	return response.SuccessWithMessage(
		c,
		fiber.StatusCreated,
		"organization created successfully",
		&CreateOrganizationResponse{
			Name:        req.GetName(),
			Description: req.GetDescription(),
		},
	)
}

// @Summary Onboard Organization
// @Description Onboard a new organization (for new users)
// @Tags Organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body tenantV1.CreateOrganizationRequest true "Organization details"
// @Success 201 {object} response.SuccessResponse{data=OnboardOrganizationResponse}
// @Router /organizations/onboarding [post]
// onboarding is called only when the user is not yet a member of any org.
func (h *OrganizationHandler) onboarding(c *fiber.Ctx) error {
	var req tenantV1.CreateOrganizationRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, ErrInvalidBody)
	}

	resp, err := h.client.OnboardOrganization(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	h.logger.Debug("Organization onboarded", resp)

	dtoResp := protoOrgResponseToDTO(resp)

	return response.SuccessWithMessage(
		c,
		fiber.StatusCreated,
		"organization onboarded successfully",
		&OnboardOrganizationResponse{
			Organization: dtoResp.Organization,
			Members:      dtoResp.Members,
		},
	)
}

// @Summary Get Organization by ID
// @Description Get organization details by ID
// @Tags Organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param organizationID path string true "Organization ID"
// @Success 200 {object} response.SuccessResponse{data=Organization}
// @Router /organizations/{organizationID} [get]
func (h *OrganizationHandler) getById(c *fiber.Ctx) error {
	req := tenantV1.TenantOrganizationIDRequest{
		OrganizationId: getOrgIDFromLocals(c),
	}

	resp, err := h.client.GetOrganizationById(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	h.logger.Debug("Organization fetched by id", resp)
	orgDto := protoOrgResponseToDTO(resp)

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"organization fetched successfully",
		orgDto,
	)
}

// @Summary Get Organization by Slug
// @Description Get organization details by its unique slug
// @Tags Organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param slug path string true "Organization Slug"
// @Success 200 {object} response.SuccessResponse{data=Organization}
// @Router /organizations/slug/{slug} [get]
func (h *OrganizationHandler) getBySlug(c *fiber.Ctx) error {
	var (
		slug string
		in   tenantV1.OrganizationSlugRequest
	)

	slug = c.Params(ParamSlug)
	if slug == "" {
		return response.Error(c, ErrInvalidSlug)
	}
	in.Slug = slug

	resp, err := h.client.GetBySlug(c.UserContext(), &in)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	h.logger.Debug("Organization fetched by slug", resp)
	orgDto := protoOrgResponseToDTO(resp)

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"organization fetched successfully",
		orgDto,
	)
}

// @Summary List Organizations
// @Description List organizations the user belongs to
// @Tags Organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page_size query int false "Page size"
// @Param page_token query string false "Page token"
// @Success 200 {object} response.SuccessResponse{data=ListOrganizationsResponse}
// @Router /organizations [get]
func (h *OrganizationHandler) list(c *fiber.Ctx) error {
	req := tenantV1.ListOrganizationsRequest{
		PageSize:  int32(c.QueryInt(QueryParamPageSize, 0)),
		PageToken: c.Query(QueryParamPageToken, ""),
	}

	resp, err := h.client.ListOrganizations(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	orgsDto := make([]*Organization, len(resp.GetOrganizations()))
	for i, org := range resp.GetOrganizations() {
		orgsDto[i] = protoOrgToDTO(org)
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"organizations fetched successfully",
		&ListOrganizationsResponse{
			Organizations: orgsDto,
			TotalCount:    resp.GetTotalCount(),
			NextPageToken: resp.GetNextPageToken(),
		},
	)
}

// @Summary Update Organization
// @Description Update an existing organization
// @Tags Organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param organizationID path string true "Organization ID"
// @Param request body tenantV1.UpdateOrganizationRequest true "Update details"
// @Success 200 {object} response.SuccessResponse{data=Organization}
// @Router /organizations/{organizationID} [patch]
func (h *OrganizationHandler) update(c *fiber.Ctx) error {
	var req tenantV1.UpdateOrganizationRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, ErrInvalidBody)
	}
	req.OrganizationId = getOrgIDFromLocals(c)

	resp, err := h.client.UpdateOrganization(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	orgDto := protoOrgResponseToDTO(resp)

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"organization updated successfully",
		orgDto,
	)
}

// @Summary Archive Organization
// @Description Archive an organization
// @Tags Organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param organizationID path string true "Organization ID"
// @Param request body tenantV1.ArchiveOrganizationRequest true "Archive request details"
// @Success 200 {object} response.SuccessResponse
// @Router /organizations/{organizationID}/archive [delete]
func (h *OrganizationHandler) archive(c *fiber.Ctx) error {
	var req tenantV1.ArchiveOrganizationRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, ErrInvalidBody)
	}
	req.OrganizationId = getOrgIDFromLocals(c)

	resp, err := h.client.ArchiveOrganization(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusOK, resp.GetMessage(), nil)
}

// @Summary Restore Organization
// @Description Restore an archived organization
// @Tags Organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param organizationID path string true "Organization ID"
// @Success 200 {object} response.SuccessResponse
// @Router /organizations/{organizationID}/restore [post]
func (h *OrganizationHandler) restore(c *fiber.Ctx) error {
	req := tenantV1.TenantOrganizationIDRequest{
		OrganizationId: getOrgIDFromLocals(c),
	}

	resp, err := h.client.RestoreOrganization(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusOK, resp.GetMessage(), nil)
}

// @Summary Delete Organization
// @Description Delete an organization completely
// @Tags Organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param organizationID path string true "Organization ID"
// @Param request body tenantV1.DeleteOrganizationRequest true "Delete request details"
// @Success 200 {object} response.SuccessResponse
// @Router /organizations/{organizationID} [delete]
func (h *OrganizationHandler) delete(c *fiber.Ctx) error {
	var req tenantV1.DeleteOrganizationRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, ErrInvalidBody)
	}
	req.OrganizationId = getOrgIDFromLocals(c)

	resp, err := h.client.DeleteOrganization(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusOK, resp.GetMessage(), nil)
}

// @Summary Transfer Organization Ownership
// @Description Transfer the ownership of the organization
// @Tags Organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param organizationID path string true "Organization ID"
// @Param request body tenantV1.TransferOwnershipRequest true "Transfer details"
// @Success 200 {object} response.SuccessResponse
// @Router /organizations/{organizationID}/transfer-ownership [post]
func (h *OrganizationHandler) transferOwnership(c *fiber.Ctx) error {
	var req tenantV1.TransferOwnershipRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, ErrInvalidBody)
	}
	req.OrganizationId = getOrgIDFromLocals(c)

	resp, err := h.client.TransferOwnership(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusOK, resp.GetMessage(), nil)
}

// @Summary List Organization Roles
// @Description Get all roles for the organization
// @Tags Organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param organizationID path string true "Organization ID"
// @Success 200 {object} response.SuccessResponse{data=[]OrganizationRole}
// @Router /organizations/{organizationID}/roles [get]
func (h *OrganizationHandler) listRoles(c *fiber.Ctx) error {
	var req tenantV1.TenantOrganizationIDRequest
	req.OrganizationId = getOrgIDFromLocals(c)

	resp, err := h.client.GetAllRoles(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	rolesDto := make([]*OrganizationRole, len(resp.GetRoles()))
	for i, r := range resp.GetRoles() {
		rolesDto[i] = protoOrgRoleToDTO(r)
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"organization roles fetched successfully",
		rolesDto,
	)
}

// @Summary List Organization Permissions
// @Description Get all permissions for the organization
// @Tags Organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param organizationID path string true "Organization ID"
// @Success 200 {object} response.SuccessResponse{data=[]OrganizationPerm}
// @Router /organizations/{organizationID}/permissions [get]
// Get all permissions
func (h *OrganizationHandler) listPermissions(c *fiber.Ctx) error {
	var req tenantV1.TenantOrganizationIDRequest
	req.OrganizationId = getOrgIDFromLocals(c)

	resp, err := h.client.GetAllPermissions(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	permsDto := make([]*OrganizationPerm, len(resp.GetPermissions()))
	for i, p := range resp.GetPermissions() {
		permsDto[i] = protoPermToDTO(p)
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"organization permissions fetched successfully",
		permsDto,
	)
}
