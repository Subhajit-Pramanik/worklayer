package apikey

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/pkg/ctxutil"
	"github.com/vyolayer/vyolayer/pkg/errors"
	"github.com/vyolayer/vyolayer/pkg/response"
	pb "github.com/vyolayer/vyolayer/proto/apikey/v1"
)

// @Summary Create API Key
// @Description Create a new API key for a project
// @Tags API Keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId query string true "Organization ID"
// @Param projectId query string true "Project ID"
// @Param request body CreateAPIKeyRequest true "API Key details"
// @Success 200 {object} response.SuccessResponse{data=CreateAPIKeyDTO}
// @Router /api-keys [post]
func (h *ApiKeyHandler) Create(c *fiber.Ctx) error {

	var (
		in             CreateAPIKeyRequest
		organizationID string
		projectID      string
		actorID        string
	)

	if err := c.BodyParser(&in); err != nil {
		return response.Error(c, errors.BadRequest("invalid request body"))
	}

	organizationID, projectID, err := checkOrgAndProject(c)
	if err != nil {
		return response.Error(c, err)
	}


	actorID, _ = ctxutil.ExtractIAMUserID(c.UserContext())
	if actorID == "" {
		return response.Error(c, errors.Unauthorized("unauthorized"))
	}

	resp, err := h.client.CreateAPIKey(c.UserContext(), &pb.CreateAPIKeyRequest{
		OrganizationId: organizationID,
		ProjectId:      projectID,
		ActorId:        actorID,
		Name:           in.Name,
		Description:    in.Description,
		Environment:    in.Environment,
		Scopes:         in.Scopes,
		// TODO: add expires at
	})
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.Success(c, resp)
}

// @Summary List API Keys
// @Description List all API keys for a project
// @Tags API Keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId query string true "Organization ID"
// @Param projectId query string true "Project ID"
// @Param page query int false "Page number"
// @Param limit query int false "Page limit"
// @Param status query string false "Status filter"
// @Success 200 {object} response.SuccessResponse{data=ListAPIKeysDTO}
// @Router /api-keys [get]
func (h *ApiKeyHandler) List(c *fiber.Ctx) error {

	var (
		organizationID string
		projectID      string
	)

	organizationID, projectID, err := checkOrgAndProject(c)
	if err != nil {
		return response.Error(c, err)
	}

	page := c.QueryInt(QueryParamPage, 1)
	limit := c.QueryInt(QueryParamLimit, 5)
	status := c.Query(QueryParamStatus)

	resp, err := h.client.ListAPIKeys(c.UserContext(), &pb.ListAPIKeysRequest{
		OrganizationId: organizationID,
		ProjectId:      projectID,
		Page:           int32(page),
		Limit:          int32(limit),
		Status:         status,
	})

	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.Success(c, resp)
}

// @Summary Get API Key
// @Description Get details of a specific API key
// @Tags API Keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId query string true "Organization ID"
// @Param projectId query string true "Project ID"
// @Param apiKeyID path string true "API Key ID"
// @Success 200 {object} response.SuccessResponse{data=GetAPIKeyDTO}
// @Router /api-keys/{apiKeyID} [get]
func (h *ApiKeyHandler) Get(c *fiber.Ctx) error {

	organizationID, projectID, err := checkOrgAndProject(c)
	if err != nil {
		return response.Error(c, err)
	}

	apiKeyID := c.Params(ParamApiKeyID)
	if apiKeyID == "" {
		return response.Error(c, errors.BadRequest("api key id is required"))
	}

	resp, err := h.client.GetAPIKey(c.UserContext(), &pb.GetAPIKeyRequest{
		OrganizationId: organizationID,
		ProjectId:      projectID,
		Id:             apiKeyID,
	})
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.Success(c, resp)
}

// @Summary Revoke API Key
// @Description Revoke a specific API key
// @Tags API Keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId query string true "Organization ID"
// @Param projectId query string true "Project ID"
// @Param apiKeyID path string true "API Key ID"
// @Success 200 {object} response.SuccessResponse{data=RevokeAPIKeyDTO}
// @Router /api-keys/{apiKeyID}/revoke [delete]
func (h *ApiKeyHandler) Revoke(c *fiber.Ctx) error {

	organizationID, projectID, err := checkOrgAndProject(c)
	if err != nil {
		return response.Error(c, err)
	}

	apiKeyID := c.Params(ParamApiKeyID)
	if apiKeyID == "" {
		return response.Error(c, errors.BadRequest("api key id is required"))
	}

	actorID, _ := ctxutil.ExtractIAMUserID(c.UserContext())
	if actorID == "" {
		return response.Error(c, errors.Unauthorized("unauthorized"))
	}

	resp, err := h.client.RevokeAPIKey(c.UserContext(), &pb.RevokeAPIKeyRequest{
		OrganizationId: organizationID,
		ProjectId:      projectID,
		Id:             apiKeyID,
		ActorId:        actorID,
	})
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.Success(c, resp)
}

// @Summary Rotate API Key
// @Description Rotate a specific API key
// @Tags API Keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId query string true "Organization ID"
// @Param projectId query string true "Project ID"
// @Param apiKeyID path string true "API Key ID"
// @Success 200 {object} response.SuccessResponse{data=RotateAPIKeyDTO}
// @Router /api-keys/{apiKeyID}/rotate [patch]
func (h *ApiKeyHandler) Rotate(c *fiber.Ctx) error {

	organizationID, projectID, err := checkOrgAndProject(c)
	if err != nil {
		return response.Error(c, err)
	}

	apiKeyID := c.Params(ParamApiKeyID)
	if apiKeyID == "" {
		return response.Error(c, errors.BadRequest("api key id is required"))
	}

	actorID, _ := ctxutil.ExtractIAMUserID(c.UserContext())
	if actorID == "" {
		return response.Error(c, errors.Unauthorized("unauthorized"))
	}

	resp, err := h.client.RotateAPIKey(c.UserContext(), &pb.RotateAPIKeyRequest{
		OrganizationId: organizationID,
		ProjectId:      projectID,
		ActorId:        actorID,
		Id:             apiKeyID,
	})
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.Success(c, resp)
}

// @Summary Validate API Key
// @Description Validate if an API key is active
// @Tags API Keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param apiKeyID path string true "API Key ID"
// @Param request body ValidateAPIKeyRequest true "Validation details"
// @Success 200 {object} response.SuccessResponse{data=ValidateAPIKeyDTO}
// @Router /api-keys/{apiKeyID}/validate [get]
func (h *ApiKeyHandler) Validate(c *fiber.Ctx) error {
	var req pb.ValidateAPIKeyRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, errors.BadRequest("invalid request body"))
	}

	resp, err := h.client.ValidateAPIKey(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.Success(c, resp)
}

// check params orgId and projectId
func checkOrgAndProject(c *fiber.Ctx) (string, string, error) {
	organizationID := c.Query(QueryParamOrganizationId)
	if organizationID == "" {
		return "", "", errors.BadRequest("organization id is required")
	}

	projectID := c.Query(QueryParamProjectId)
	if projectID == "" {
		return "", "", errors.BadRequest("project id is required")
	}

	return organizationID, projectID, nil
}
