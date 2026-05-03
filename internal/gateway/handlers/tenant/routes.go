package tenant

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/internal/gateway/middleware"
)

func (h *OrganizationHandler) RegisterRoutes(router fiber.Router) {
	grpcCtxMiddleware := middleware.NewGrpcCtxMiddleware(tenantGRPCTimeout).Handler()

	org := router.Group("/organizations")
	org.Use(grpcCtxMiddleware)
	org.Use(middleware.IamJWTVerify(h.authClient))

	org.
		Post("/onboarding", h.onboarding).
		Post("/", h.create).
		Get("/", h.list)

	org.Get("/slug/:"+ParamSlug, h.getBySlug)

	// Invitation routes that don't require an org context (accept uses token, pending is user-scoped)
	org.Post("/invitations/accept", h.invitationHandler.acceptInvitation)
	org.Get("/invitations/pending", h.invitationHandler.getPendingByUser)

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

	// Invitations
	orgGroup.
		Post("/invitations", h.invitationHandler.createInvitation).
		Get("/invitations", h.invitationHandler.listInvitations).
		Get("/invitations/pending", h.invitationHandler.getPendingByOrgID).
		Delete("/invitations/:"+ParamInvitationID, h.invitationHandler.cancelInvitation)

	orgMemberGroup := orgGroup.Group("/members")
	orgMemberGroup.
		Get("/", h.memberHandler.listMembers).
		Get("/me", h.memberHandler.getCurrentMember).
		Get("/:"+ParamMemberID, h.memberHandler.getMemberByID).
		Delete("/:"+ParamMemberID, h.memberHandler.removeMember)

	h.logger.Info("Organization routes registered", "")
}

func (h *ProjectHandler) RegisterRoutes(router fiber.Router) {
	grpcCtxMiddleware := middleware.NewGrpcCtxMiddleware(tenantGRPCTimeout)

	projects := router.Group("/organizations/:" + ParamOrganizationID + "/projects")
	projects.Use(
		grpcCtxMiddleware.Handler(),
		middleware.IamJWTVerify(h.authClient),
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
