package wire

import (
	accounthandler "github.com/vyolayer/vyolayer/internal/gateway/handlers/account"
	apikeyhandler "github.com/vyolayer/vyolayer/internal/gateway/handlers/apikey"
	consolehandler "github.com/vyolayer/vyolayer/internal/gateway/handlers/console"
	healthhandler "github.com/vyolayer/vyolayer/internal/gateway/handlers/health"
	iamhandler "github.com/vyolayer/vyolayer/internal/gateway/handlers/iam"
	tenanthandler "github.com/vyolayer/vyolayer/internal/gateway/handlers/tenant"
	"github.com/vyolayer/vyolayer/internal/gateway/server"
	"github.com/vyolayer/vyolayer/internal/gateway/service"
	"github.com/vyolayer/vyolayer/pkg/jwt"
	"github.com/vyolayer/vyolayer/pkg/logger"
)

func NewRegistrars(
	logger *logger.AppLogger,
	clients *Clients,
	cookieSrv *service.AccountTokenService,
	accountJWT jwt.AccountJWT,
	iamCookieSrv *service.IAMCookieService,
) []server.RouteRegistrar {
	return []server.RouteRegistrar{
		healthhandler.NewHealthHandler(),

		// Account routes
		accounthandler.NewAccountHandler(
			clients.AccountClient,
			cookieSrv,
			accountJWT,
			logger,
		),

		// IAM routes
		iamhandler.NewIAMAuthGatewayHandler(
			clients.IamAuthClient,
			clients.IamUserClient,
			iamCookieSrv,
			clients.IamAuthClient,
			logger,
		),

		// Tenant routes
		// Tenant Organization routes
		tenanthandler.NewOrganizationHandler(
			logger,
			clients.TenantOrganizationClient,
			clients.IamAuthClient,
		),

		// Tenant Organization Member routes
		tenanthandler.NewOrganizationMemberHandler(
			logger,
			clients.TenantOrganizationMemClient,
			clients.IamAuthClient,
		),

		// Tenant Organization Invitation routes
		tenanthandler.NewOrganizationInvitationHandler(
			logger,
			clients.TenantOrganizationInvClient,
			clients.IamAuthClient,
		),

		// Tenant Project & Project Member routes
		tenanthandler.NewProjectHandler(
			logger,
			clients.TenantProjectClient,
			clients.IamAuthClient,
		),

		consolehandler.NewProjectServiceHandler(
			logger,
			clients.ConsoleProjectServiceManifestClient,
			clients.IamAuthClient,
		),

		apikeyhandler.NewHandler(
			logger,
			clients.ApikeyServiceClient,
			clients.IamAuthClient,
		),
	}
}
