package tenant

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/internal/gateway/middleware"
)

func (h *OrganizationHandler) RegisterRoutes(router fiber.Router) {
	grpcCtxMiddleware := middleware.NewGrpcCtxMiddleware(tenantGRPCTimeout).Handler()

	org := router.Group("/organizations")
	org.Use(grpcCtxMiddleware)
	org.Use(middleware.IamJWTVerify(h.iamJWT))

	org.
		Post("/onboarding", h.onboarding).
		Post("/", h.create).
		Get("/", h.list)

	org.Get("/slug/:"+ParamSlug, h.getBySlug)

	// All routes below require a valid organizationID in the path
	orgGroup := org.Group("/:"+ParamOrganizationID, middleware.ValidateOrganizationID())

	// Organization lifecycle
	orgGroup.
		Get("/", h.getById).
		Patch("/", h.update).
		Delete("/", h.delete).
		Delete("/archive", h.archive).
		Post("/restore", h.restore).
		Post("/transfer-ownership", h.transferOwnership)

	// Roles and permissions
	orgGroup.Get("/roles", h.listRoles).
		Get("/permissions", h.listPermissions)

	h.logger.Info("Organization routes registered", "")
}

func (h *OrganizationInvitationHandler) RegisterRoutes(router fiber.Router) {
	grpcCtxMiddleware := middleware.NewGrpcCtxMiddleware(tenantGRPCTimeout).Handler()

	org := router.Group("/organizations")
	org.Use(grpcCtxMiddleware, middleware.IamJWTVerify(h.iamJWT))

	// Invitation routes that don't require an org context (accept uses token, pending is user-scoped)
	org.Post("/invitations/accept", h.acceptInvitation)
	org.Get("/invitations/pending", h.getPendingByUser)

	// All routes below require a valid organizationID in the path
	orgGroup := org.Group("/:"+ParamOrganizationID, middleware.ValidateOrganizationID())

	// Invitations
	orgGroup.
		Post("/invitations", h.createInvitation).
		Get("/invitations", h.listInvitations).
		Get("/invitations/pending", h.getPendingByOrgID).
		Delete("/invitations/:"+ParamInvitationID, h.cancelInvitation)

	h.logger.Info("Organization invitation routes registered", "")
}

func (h *OrganizationMemberHandler) RegisterRoutes(router fiber.Router) {
	grpcCtxMiddleware := middleware.NewGrpcCtxMiddleware(tenantGRPCTimeout).Handler()

	orgMemberGroup := router.Group("/organizations/:" + ParamOrganizationID + "/members")
	orgMemberGroup.Use(
		grpcCtxMiddleware,
		middleware.IamJWTVerify(h.iamJWT),
		middleware.ValidateOrganizationID(),
	)

	// Members
	orgMemberGroup.
		Get("/", h.listMembers).
		Get("/me", h.getCurrentMember).
		Get("/:"+ParamMemberID, h.getMemberByID).
		Delete("/:"+ParamMemberID, h.removeMember)

	h.logger.Info("Organization member routes registered", "")
}

func (h *ProjectHandler) RegisterRoutes(router fiber.Router) {
	grpcCtxMiddleware := middleware.NewGrpcCtxMiddleware(tenantGRPCTimeout)

	projects := router.Group("/organizations/:" + ParamOrganizationID + "/projects")
	projects.Use(
		grpcCtxMiddleware.Handler(),
		middleware.IamJWTVerify(h.iamJWT),
		middleware.ValidateOrganizationID(),
	)

	projects.Get("/", h.listProjects)
	projects.Post("/", h.createProject)

	project := projects.Group("/:"+ParamProjectID, middleware.ValidateProjectID())
	project.Get("/", h.getProject)
	project.Patch("/", h.updateProject)
	project.Delete("/", h.deleteProject)

	members := project.Group("/members")
	members.Get("/", h.listMembers)
	members.Post("/", h.addMember)
	members.Get("/me", h.getCurrentMember)
	// members.Delete("/leave", h.leaveProject)
	members.Get("/:"+ParamMemberID, h.getMember)
	members.Post("/:"+ParamMemberID+"/role", h.changeMemberRole)
	members.Delete("/:"+ParamMemberID, h.removeMember)

	h.logger.Info("Project routes registered", "")
}
