package tenant

import (
	"github.com/vyolayer/vyolayer/pkg/logger"
	iAMV1 "github.com/vyolayer/vyolayer/proto/iam/v1"
	tenantV1 "github.com/vyolayer/vyolayer/proto/tenant/v1"
)

type OrganizationHandler struct {
	logger     *logger.AppLogger
	client     tenantV1.OrganizationServiceClient
	authClient iAMV1.AuthServiceClient
}

func NewOrganizationHandler(
	logger *logger.AppLogger,
	client tenantV1.OrganizationServiceClient,
	authClient iAMV1.AuthServiceClient,
) *OrganizationHandler {
	return &OrganizationHandler{
		logger:     logger.WithContext("Org Handler"),
		client:     client,
		authClient: authClient,
	}
}

type OrganizationInvitationHandler struct {
	logger     *logger.AppLogger
	client     tenantV1.OrganizationInvitationServiceClient
	authClient iAMV1.AuthServiceClient
}

func NewOrganizationInvitationHandler(
	logger *logger.AppLogger,
	client tenantV1.OrganizationInvitationServiceClient,
	authClient iAMV1.AuthServiceClient,
) *OrganizationInvitationHandler {
	return &OrganizationInvitationHandler{
		logger:     logger.WithContext("Org Invitation Handler"),
		client:     client,
		authClient: authClient,
	}
}

type OrganizationMemberHandler struct {
	logger     *logger.AppLogger
	client     tenantV1.OrganizationMemberServiceClient
	authClient iAMV1.AuthServiceClient
}

func NewOrganizationMemberHandler(
	logger *logger.AppLogger,
	client tenantV1.OrganizationMemberServiceClient,
	authClient iAMV1.AuthServiceClient,
) *OrganizationMemberHandler {
	return &OrganizationMemberHandler{
		logger:     logger.WithContext("Org Member Handler"),
		client:     client,
		authClient: authClient,
	}
}

type ProjectHandler struct {
	logger     *logger.AppLogger
	client     tenantV1.ProjectServiceClient
	authClient iAMV1.AuthServiceClient
}

func NewProjectHandler(
	logger *logger.AppLogger,
	client tenantV1.ProjectServiceClient,
	authClient iAMV1.AuthServiceClient,
) *ProjectHandler {
	return &ProjectHandler{
		logger:     logger.WithContext("Project Handler"),
		client:     client,
		authClient: authClient,
	}
}
