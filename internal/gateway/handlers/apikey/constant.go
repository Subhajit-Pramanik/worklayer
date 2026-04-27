package apikey

import "time"

const (
	QueryParamOrganizationId = "orgId"
	QueryParamProjectId      = "projectId"
	QueryParamPage           = "page"
	QueryParamLimit          = "limit"
	QueryParamStatus         = "status"
	ParamApiKeyID            = "apiKeyID"
	grpcTimeout              = 10 * time.Second
)
