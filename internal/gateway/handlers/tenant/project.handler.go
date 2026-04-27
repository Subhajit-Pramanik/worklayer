package tenant

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/pkg/errors"
	"github.com/vyolayer/vyolayer/pkg/response"
	tenantV1 "github.com/vyolayer/vyolayer/proto/tenant/v1"
)

// ─── Project CRUD ─────────────────────────────────────────────────────────────

// @Summary Create Project
// @Description Create a new project within an organization
// @Tags Projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param organizationID path string true "Organization ID"
// @Param request body tenantV1.CreateProjectRequest true "Project details"
// @Success 201 {object} response.SuccessResponse{data=Project}
// @Router /organizations/{organizationID}/projects [post]
func (h *ProjectHandler) createProject(c *fiber.Ctx) error {
	var req tenantV1.CreateProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, ErrInvalidBody)
	}
	req.OrganizationId = getOrgIDFromLocals(c)

	resp, err := h.client.CreateProject(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusCreated,
		"project created successfully",
		protoProjectResponseToDTO(resp),
	)
}

// @Summary Get Project
// @Description Get details of a specific project
// @Tags Projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param organizationID path string true "Organization ID"
// @Param projectID path string true "Project ID"
// @Success 200 {object} response.SuccessResponse{data=Project}
// @Router /organizations/{organizationID}/projects/{projectID} [get]
func (h *ProjectHandler) getProject(c *fiber.Ctx) error {
	req := tenantV1.GetProjectRequest{
		OrganizationId: getOrgIDFromLocals(c),
		ProjectId:      getProjectIDFromLocals(c),
	}

	resp, err := h.client.GetProject(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"project fetched successfully",
		protoProjectResponseToDTO(resp),
	)
}

// @Summary List Projects
// @Description List projects within an organization
// @Tags Projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param organizationID path string true "Organization ID"
// @Param page_size query int false "Page size"
// @Param page_token query string false "Page token"
// @Success 200 {object} response.SuccessResponse{data=ListProjectsResponse}
// @Router /organizations/{organizationID}/projects [get]
func (h *ProjectHandler) listProjects(c *fiber.Ctx) error {
	req := tenantV1.ListProjectsRequest{
		OrganizationId: getOrgIDFromLocals(c),
		PageSize:       int32(c.QueryInt(QueryParamPageSize, 0)),
		PageToken:      c.Query(QueryParamPageToken, ""),
	}

	resp, err := h.client.ListProjects(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	projects := make([]*Project, len(resp.GetProjects()))
	for i, p := range resp.GetProjects() {
		projects[i] = protoProjectToDTO(p)
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"projects fetched successfully",
		&ListProjectsResponse{
			Projects:      projects,
			TotalCount:    resp.GetTotalCount(),
			NextPageToken: resp.GetNextPageToken(),
		},
	)
}

// @Summary Update Project
// @Description Update project details
// @Tags Projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param organizationID path string true "Organization ID"
// @Param projectID path string true "Project ID"
// @Param request body tenantV1.UpdateProjectRequest true "Update details"
// @Success 200 {object} response.SuccessResponse{data=Project}
// @Router /organizations/{organizationID}/projects/{projectID} [patch]
func (h *ProjectHandler) updateProject(c *fiber.Ctx) error {
	var req tenantV1.UpdateProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, ErrInvalidBody)
	}
	req.OrganizationId = getOrgIDFromLocals(c)
	req.ProjectId = getProjectIDFromLocals(c)

	resp, err := h.client.UpdateProject(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"project updated successfully",
		protoProjectToDTO(resp.GetProject()),
	)
}

// @Summary Delete Project
// @Description Delete a project
// @Tags Projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param organizationID path string true "Organization ID"
// @Param projectID path string true "Project ID"
// @Param request body tenantV1.DeleteProjectRequest true "Delete details"
// @Success 200 {object} response.SuccessResponse
// @Router /organizations/{organizationID}/projects/{projectID} [delete]
func (h *ProjectHandler) deleteProject(c *fiber.Ctx) error {
	var req tenantV1.DeleteProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, ErrInvalidBody)
	}
	req.OrganizationId = getOrgIDFromLocals(c)
	req.ProjectId = getProjectIDFromLocals(c)

	resp, err := h.client.DeleteProject(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusOK, resp.GetMessage(), nil)
}

// ─── Project Member operations ────────────────────────────────────────────────

// @Summary List Project Members
// @Description List members of a project
// @Tags Project Members
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param organizationID path string true "Organization ID"
// @Param projectID path string true "Project ID"
// @Param page_size query int false "Page size"
// @Param page_token query string false "Page token"
// @Success 200 {object} response.SuccessResponse{data=ListProjectMembersResponse}
// @Router /organizations/{organizationID}/projects/{projectID}/members [get]
func (h *ProjectHandler) listMembers(c *fiber.Ctx) error {
	req := tenantV1.ListProjectMembersRequest{
		OrganizationId: getOrgIDFromLocals(c),
		ProjectId:      getProjectIDFromLocals(c),
		PageSize:       int32(c.QueryInt(QueryParamPageSize, 0)),
		PageToken:      c.Query(QueryParamPageToken, ""),
	}

	resp, err := h.client.ListMembers(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	members := make([]*ProjectMember, len(resp.GetMembers()))
	for i, m := range resp.GetMembers() {
		members[i] = protoProjectMemberToDTO(m)
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"members fetched successfully",
		&ListProjectMembersResponse{
			Members:       members,
			TotalCount:    resp.GetTotalCount(),
			NextPageToken: resp.GetNextPageToken(),
		},
	)
}

// @Summary Get Project Member
// @Description Get details of a project member
// @Tags Project Members
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param organizationID path string true "Organization ID"
// @Param projectID path string true "Project ID"
// @Param memberID path string true "Member ID"
// @Success 200 {object} response.SuccessResponse{data=ProjectMember}
// @Router /organizations/{organizationID}/projects/{projectID}/members/{memberID} [get]
func (h *ProjectHandler) getMember(c *fiber.Ctx) error {
	req := tenantV1.GetProjectMemberRequest{
		OrganizationId: getOrgIDFromLocals(c),
		ProjectId:      getProjectIDFromLocals(c),
		MemberId:       c.Params(ParamMemberID),
	}

	resp, err := h.client.GetMember(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"member fetched successfully",
		protoProjectMemberToDTO(resp.GetMember()),
	)
}

// @Summary Get Current Project Member
// @Description Get current user's membership details for the project
// @Tags Project Members
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param organizationID path string true "Organization ID"
// @Param projectID path string true "Project ID"
// @Success 200 {object} response.SuccessResponse{data=ProjectMember}
// @Router /organizations/{organizationID}/projects/{projectID}/members/me [get]
func (h *ProjectHandler) getCurrentMember(c *fiber.Ctx) error {
	req := tenantV1.ListProjectMembersRequest{
		OrganizationId: getOrgIDFromLocals(c),
		ProjectId:      getProjectIDFromLocals(c),
	}

	resp, err := h.client.GetCurrentMember(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"current member fetched successfully",
		protoProjectMemberToDTO(resp.GetMember()),
	)
}

// @Summary Add Project Member
// @Description Add a member to a project
// @Tags Project Members
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param organizationID path string true "Organization ID"
// @Param projectID path string true "Project ID"
// @Param request body tenantV1.AddProjectMemberRequest true "Member details"
// @Success 201 {object} response.SuccessResponse{data=ProjectMember}
// @Router /organizations/{organizationID}/projects/{projectID}/members [post]
func (h *ProjectHandler) addMember(c *fiber.Ctx) error {
	var req tenantV1.AddProjectMemberRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, ErrInvalidBody)
	}
	req.OrganizationId = getOrgIDFromLocals(c)
	req.ProjectId = getProjectIDFromLocals(c)

	resp, err := h.client.AddMember(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusCreated,
		"member added successfully",
		protoProjectMemberToDTO(resp.GetMember()),
	)
}

// @Summary Change Project Member Role
// @Description Change the role of a project member
// @Tags Project Members
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param organizationID path string true "Organization ID"
// @Param projectID path string true "Project ID"
// @Param memberID path string true "Member ID"
// @Param request body tenantV1.ChangeProjectMemberRoleRequest true "Role details"
// @Success 200 {object} response.SuccessResponse
// @Router /organizations/{organizationID}/projects/{projectID}/members/{memberID}/role [post]
func (h *ProjectHandler) changeMemberRole(c *fiber.Ctx) error {
	var req tenantV1.ChangeProjectMemberRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, ErrInvalidBody)
	}
	req.OrganizationId = getOrgIDFromLocals(c)
	req.ProjectId = getProjectIDFromLocals(c)
	req.MemberId = c.Params(ParamMemberID)

	resp, err := h.client.ChangeMemberRole(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusOK, resp.GetMessage(), nil)
}

// @Summary Remove Project Member
// @Description Remove a member from the project
// @Tags Project Members
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param organizationID path string true "Organization ID"
// @Param projectID path string true "Project ID"
// @Param memberID path string true "Member ID"
// @Success 200 {object} response.SuccessResponse
// @Router /organizations/{organizationID}/projects/{projectID}/members/{memberID} [delete]
func (h *ProjectHandler) removeMember(c *fiber.Ctx) error {
	req := tenantV1.RemoveProjectMemberRequest{
		OrganizationId: getOrgIDFromLocals(c),
		ProjectId:      getProjectIDFromLocals(c),
		MemberId:       c.Params(ParamMemberID),
	}

	resp, err := h.client.RemoveMember(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusOK, resp.GetMessage(), nil)
}

// @Summary Leave Project
// @Description Leave the project
// @Tags Project Members
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param organizationID path string true "Organization ID"
// @Param projectID path string true "Project ID"
// @Success 200 {object} response.SuccessResponse
// @Router /organizations/{organizationID}/projects/{projectID}/members/leave [delete]
func (h *ProjectHandler) leaveProject(c *fiber.Ctx) error {
	req := tenantV1.ProjectIdRequest{
		OrganizationId: getOrgIDFromLocals(c),
		ProjectId:      getProjectIDFromLocals(c),
	}

	resp, err := h.client.LeaveProject(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusOK, resp.GetMessage(), nil)
}
