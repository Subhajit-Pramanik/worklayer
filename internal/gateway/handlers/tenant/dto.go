package tenant

// ── Organization ──────────────────────────────────────────────────────────────

type Organization struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Description  string `json:"description"`
	IsActive     bool   `json:"is_active"`
	OwnerID      string `json:"owner_id"`
	MaxMembers   uint32 `json:"max_members"`
	MaxProjects  uint32 `json:"max_projects"`
	ProjectCount uint32 `json:"project_count"`
	MemberCount  uint32 `json:"member_count"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type CreateOrganizationResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type OnboardOrganizationResponse struct {
	Organization *Organization         `json:"organization"`
	Members      []*OrganizationMember `json:"members,omitempty"`
}

type OrganizationDetailResponse struct {
	Organization *Organization         `json:"organization"`
	Members      []*OrganizationMember `json:"members,omitempty"`
}

type ListOrganizationsResponse struct {
	Organizations []*Organization `json:"organizations"`
	TotalCount    int32           `json:"total_count"`
	NextPageToken string          `json:"next_page_token"`
}

// ── Organization RBAC ─────────────────────────────────────────────────────────

type OrganizationRole struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	IsSystemRole bool   `json:"is_system_role"`
	IsDefault    bool   `json:"is_default"`
}

type OrganizationPerm struct {
	ID           string `json:"id"`
	Resource     string `json:"resource"`
	Action       string `json:"action"`
	Code         string `json:"code"`
	Group        string `json:"group"`
	IsSystemPerm bool   `json:"is_system_perm"`
}

// ── Organization Member ───────────────────────────────────────────────────────

type OrganizationMember struct {
	ID             string   `json:"id"`
	OrganizationID string   `json:"organization_id"`
	UserID         string   `json:"user_id"`
	FullName       string   `json:"full_name"`
	Email          string   `json:"email"`
	Roles          []string `json:"roles"`
	Status         string   `json:"status"`
	JoinedAt       string   `json:"joined_at"`
	InvitedAt      string   `json:"invited_at"`
	InvitedBy      string   `json:"invited_by"`
	DeactivatedBy  string   `json:"deactivated_by"`
	DeactivatedAt  string   `json:"deactivated_at"`
}

type OrganizationMemberWithRBACResponse struct {
	OrganizationMember
	Roles []string `json:"roles"`
	Perms []string `json:"perms"`
}

type ListOrganizationMembersResponse struct {
	Members    []*OrganizationMember `json:"members"`
	TotalCount int32                 `json:"total_count"`
}

// ── Organization Invitation ───────────────────────────────────────────────────

type OrganizationInvitation struct {
	ID             string   `json:"id"`
	OrganizationID string   `json:"organization_id"`
	Email          string   `json:"email"`
	RoleIDs        []string `json:"role_ids"`
	InvitedBy      string   `json:"invited_by"`
	InvitedAt      string   `json:"invited_at"`
	IsAccepted     bool     `json:"is_accepted"`
	AcceptedAt     string   `json:"accepted_at"`
	ExpiredAt      string   `json:"expired_at"`
	IsPending      bool     `json:"is_pending"`
}

type InvitedBy struct {
	MemberID string `json:"member_id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type OrganizationInvitationForOrg struct {
	Invitation *OrganizationInvitation `json:"invitation"`
	InvitedBy  *InvitedBy              `json:"invited_by"`
}

type ListOrganizationInvitationsResponse struct {
	Invitations []*OrganizationInvitation `json:"invitations"`
}

type ListOrganizationInvitationsForOrgResponse struct {
	Invitations []*OrganizationInvitationForOrg `json:"invitations"`
}

// ── Project ───────────────────────────────────────────────────────────────────

type Project struct {
	ID             string `json:"id,omitempty"`
	OrganizationID string `json:"organization_id,omitempty"`
	Name           string `json:"name,omitempty"`
	Slug           string `json:"slug,omitempty"`
	Description    string `json:"description,omitempty"`
	IsActive       bool   `json:"is_active,omitempty"`
	CreatedBy      string `json:"created_by,omitempty"`
	MaxAPIKeys     uint32 `json:"max_api_keys,omitempty"`
	MaxMembers     uint32 `json:"max_members,omitempty"`
	MemberCount    uint32 `json:"member_count,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
}

type ProjectMember struct {
	ID        string  `json:"id,omitempty"`
	UserID    string  `json:"user_id,omitempty"`
	Email     string  `json:"email,omitempty"`
	FullName  string  `json:"full_name,omitempty"`
	Role      string  `json:"role,omitempty"`
	IsActive  bool    `json:"is_active,omitempty"`
	JoinedAt  string  `json:"joined_at,omitempty"`
	RemovedAt *string `json:"removed_at,omitempty"`
}

type ListProjectsResponse struct {
	Projects      []*Project `json:"projects,omitempty"`
	TotalCount    int32      `json:"total_count,omitempty"`
	NextPageToken string     `json:"next_page_token,omitempty"`
}

type ListProjectMembersResponse struct {
	Members       []*ProjectMember `json:"members,omitempty"`
	TotalCount    int32            `json:"total_count,omitempty"`
	NextPageToken string           `json:"next_page_token,omitempty"`
}

type ProjectResponse struct {
	Project *Project         `json:"project,omitempty"`
	Members []*ProjectMember `json:"members,omitempty"`
}
