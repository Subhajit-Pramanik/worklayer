package tenant

import (
	"github.com/vyolayer/vyolayer/pkg/jwt"
	"github.com/vyolayer/vyolayer/pkg/logger"
	tenantV1 "github.com/vyolayer/vyolayer/proto/tenant/v1"
)

type OrganizationHandler struct {
	logger *logger.AppLogger
	client tenantV1.OrganizationServiceClient
	iamJWT jwt.IamJWT
}

func NewOrganizationHandler(
	logger *logger.AppLogger,
	client tenantV1.OrganizationServiceClient,
	iamJWT jwt.IamJWT,
) *OrganizationHandler {
	return &OrganizationHandler{
		logger: logger.WithContext("Org Handler"),
		client: client,
		iamJWT: iamJWT,
	}
}

type OrganizationInvitationHandler struct {
	logger *logger.AppLogger
	client tenantV1.OrganizationInvitationServiceClient
	iamJWT jwt.IamJWT
}

func NewOrganizationInvitationHandler(
	logger *logger.AppLogger,
	client tenantV1.OrganizationInvitationServiceClient,
	iamJWT jwt.IamJWT,
) *OrganizationInvitationHandler {
	return &OrganizationInvitationHandler{
		logger: logger.WithContext("Org Invitation Handler"),
		client: client,
		iamJWT: iamJWT,
	}
}

type OrganizationMemberHandler struct {
	logger *logger.AppLogger
	client tenantV1.OrganizationMemberServiceClient
	iamJWT jwt.IamJWT
}

func NewOrganizationMemberHandler(
	logger *logger.AppLogger,
	client tenantV1.OrganizationMemberServiceClient,
	iamJWT jwt.IamJWT,
) *OrganizationMemberHandler {
	return &OrganizationMemberHandler{
		logger: logger.WithContext("Org Member Handler"),
		client: client,
		iamJWT: iamJWT,
	}
}

type ProjectHandler struct {
	logger *logger.AppLogger
	client tenantV1.ProjectServiceClient
	iamJWT jwt.IamJWT
}

func NewProjectHandler(
	logger *logger.AppLogger,
	client tenantV1.ProjectServiceClient,
	iamJWT jwt.IamJWT,
) *ProjectHandler {
	return &ProjectHandler{
		logger: logger.WithContext("Project Handler"),
		client: client,
		iamJWT: iamJWT,
	}
}
