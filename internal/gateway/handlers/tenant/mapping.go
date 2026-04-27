package tenant

import (
	tenantV1 "github.com/vyolayer/vyolayer/proto/tenant/v1"
)

// Proto -> DTO (Organization)
func protoOrgToDTO(org *tenantV1.Organization) *Organization {
	if org == nil {
		return nil
	}

	return &Organization{
		ID:           org.GetId(),
		Name:         org.GetName(),
		Slug:         org.GetSlug(),
		Description:  org.GetDescription(),
		IsActive:     org.GetIsActive(),
		OwnerID:      org.GetOwnerId(),
		MaxMembers:   org.GetMaxMembers(),
		MaxProjects:  org.GetMaxProjects(),
		ProjectCount: org.GetProjectCount(),
		MemberCount:  org.GetMemberCount(),
		CreatedAt:    org.GetCreatedAt(),
		UpdatedAt:    org.GetUpdatedAt(),
	}
}

// Proto -> DTO (OrganizationResponse -> OrganizationDetailResponse)
func protoOrgResponseToDTO(resp *tenantV1.OrganizationResponse) *OrganizationDetailResponse {
	if resp == nil {
		return nil
	}

	members := make([]*OrganizationMember, len(resp.GetMembers()))
	for i, m := range resp.GetMembers() {
		members[i] = protoMemberToDTO(m)
	}

	return &OrganizationDetailResponse{
		Organization: protoOrgToDTO(resp.GetOrganization()),
		Members:      members,
	}
}

// Proto -> DTO (Organization member)
func protoMemberToDTO(m *tenantV1.OrganizationMember) *OrganizationMember {
	if m == nil {
		return nil
	}

	roleNames := make([]string, len(m.GetRoles()))
	for i, r := range m.GetRoles() {
		roleNames[i] = r.GetName()
	}

	return &OrganizationMember{
		ID:            m.GetId(),
		UserID:        m.GetUserId(),
		FullName:      m.GetFullName(),
		Email:         m.GetEmail(),
		Status:        m.GetStatus(),
		JoinedAt:      m.GetJoinedAt(),
		InvitedAt:     m.GetInvitedAt(),
		InvitedBy:     m.GetInvitedBy(),
		DeactivatedBy: m.GetDeactivatedBy(),
		DeactivatedAt: m.GetDeactivatedAt(),
		Roles:         roleNames,
	}
}

// Proto -> DTO (Organization role)
func protoOrgRoleToDTO(r *tenantV1.OrganizationRole) *OrganizationRole {
	if r == nil {
		return nil
	}

	return &OrganizationRole{
		ID:           r.GetId(),
		Name:         r.GetName(),
		Description:  r.GetDescription(),
		IsSystemRole: r.GetIsSystemRole(),
		IsDefault:    r.GetIsDefault(),
	}
}

// Proto -> DTO (Organization permission)
func protoPermToDTO(p *tenantV1.OrganizationPermission) *OrganizationPerm {
	if p == nil {
		return nil
	}

	return &OrganizationPerm{
		ID:           p.GetId(),
		Resource:     p.GetResource(),
		Action:       p.GetAction(),
		Code:         p.GetCode(),
		Group:        p.GetGroup(),
		IsSystemPerm: p.GetIsSystem(),
	}
}

// Proto -> DTO (Organization invitation)
func protoInvitationToDTO(inv *tenantV1.OrganizationMemberInvitation) *OrganizationInvitation {
	if inv == nil {
		return nil
	}

	return &OrganizationInvitation{
		ID:             inv.GetId(),
		OrganizationID: inv.GetOrganizationId(),
		Email:          inv.GetEmail(),
		RoleIDs:        inv.GetRoleIds(),
		InvitedBy:      inv.GetInvitedBy(),
		InvitedAt:      inv.GetInvitedAt(),
		IsAccepted:     inv.GetIsAccepted(),
		AcceptedAt:     inv.GetAcceptedAt(),
		ExpiredAt:      inv.GetExpiredAt(),
		IsPending:      inv.GetIsPending(),
	}
}

// Proto -> DTO (Organization invitation for org)
func protoInvitationForOrgToDTO(inv *tenantV1.OrganizationMemberInvitationForOrg) *OrganizationInvitationForOrg {
	if inv == nil {
		return nil
	}

	invDto := protoInvitationToDTO(inv.GetInvitation())
	invByDto := &InvitedBy{
		MemberID: inv.GetInvitedBy().GetMemberId(),
		FullName: inv.GetInvitedBy().GetFullName(),
		Email:    inv.GetInvitedBy().GetEmail(),
	}

	return &OrganizationInvitationForOrg{
		Invitation: invDto,
		InvitedBy:  invByDto,
	}
}

// Proto -> DTO (Project)
func protoProjectToDTO(p *tenantV1.Project) *Project {
	if p == nil {
		return nil
	}

	return &Project{
		ID:             p.GetId(),
		OrganizationID: p.GetOrganizationId(),
		Name:           p.GetName(),
		Slug:           p.GetSlug(),
		Description:    p.GetDescription(),
		IsActive:       p.GetIsActive(),
		CreatedBy:      p.GetCreatedBy(),
		MaxAPIKeys:     p.GetMaxApiKeys(),
		MaxMembers:     p.GetMaxMembers(),
		MemberCount:    p.GetMemberCount(),
		CreatedAt:      p.GetCreatedAt(),
	}
}

// Proto -> DTO (ProjectResponse)
func protoProjectResponseToDTO(resp *tenantV1.ProjectResponse) *ProjectResponse {
	if resp == nil {
		return nil
	}

	return &ProjectResponse{
		Project: protoProjectToDTO(resp.GetProject()),
	}
}

// Proto -> DTO (ProjectMember)
func protoProjectMemberToDTO(m *tenantV1.ProjectMember) *ProjectMember {
	if m == nil {
		return nil
	}

	return &ProjectMember{
		ID:        m.GetId(),
		UserID:    m.GetUserId(),
		Email:     m.GetEmail(),
		FullName:  m.GetFullName(),
		Role:      m.GetRole(),
		IsActive:  m.GetIsActive(),
		JoinedAt:  m.GetJoinedAt(),
		RemovedAt: m.RemovedAt,
	}
}
