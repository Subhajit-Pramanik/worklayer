package tenant

import (
	"time"

	"github.com/vyolayer/vyolayer/pkg/errors"
)

const (
	tenantGRPCTimeout   = 10 * time.Second
	ParamOrganizationID = "organizationID"
	ParamProjectID      = "projectID"
	ParamMemberID       = "memberID"
	ParamSlug           = "slug"
	ParamInvitationID   = "invitationID"

	QueryParamPageSize  = "page_size"
	QueryParamPageToken = "page_token"
	QueryParamToken     = "token"
)

var (
	ErrInvalidBody  = errors.BadRequest("invalid request body")
	ErrInvalidOrgID = errors.BadRequest("invalid organization id")
	ErrInvalidSlug  = errors.BadRequest("invalid slug")
)
