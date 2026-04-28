package apikey

import "time"

const (
	QueryParamOrganizationId = "org-id"
	QueryParamProjectId      = "project-id"
	QueryParamPage           = "page"
	QueryParamLimit          = "limit"
	QueryParamStatus         = "status"
	ParamApiKeyID            = "apiKeyID"
	grpcTimeout              = 10 * time.Second

	defaultPage      = 1
	defaultPageLimit = 5
)
