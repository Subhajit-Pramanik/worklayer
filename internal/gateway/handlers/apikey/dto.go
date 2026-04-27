package apikey

// APIKeyDTO is a swagger-friendly representation of the APIKey proto message.
// timestamppb.Timestamp fields are represented as RFC3339 strings.
type APIKeyDTO struct {
	Id             string   `json:"id,omitempty"`
	OrganizationId string   `json:"organization_id,omitempty"`
	ProjectId      string   `json:"project_id,omitempty"`
	Name           string   `json:"name,omitempty"`
	Description    string   `json:"description,omitempty"`
	Prefix         string   `json:"prefix,omitempty"`
	Environment    string   `json:"environment,omitempty"`
	Status         string   `json:"status,omitempty"`
	Scopes         []string `json:"scopes,omitempty"`
	CreatedBy      string   `json:"created_by,omitempty"`
	LastUsedAt     string   `json:"last_used_at,omitempty"`
	LastUsedIp     string   `json:"last_used_ip,omitempty"`
	LastUsedUa     string   `json:"last_used_ua,omitempty"`
	ExpiresAt      string   `json:"expires_at,omitempty"`
	RevokedBy      string   `json:"revoked_by,omitempty"`
	RevokedAt      string   `json:"revoked_at,omitempty"`
	CreatedAt      string   `json:"created_at,omitempty"`
	UpdatedAt      string   `json:"updated_at,omitempty"`
}

// CreateAPIKeyRequest is a request body for creating an API key.
type CreateAPIKeyRequest struct {
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Environment string   `json:"environment,omitempty"`
	Scopes      []string `json:"scopes,omitempty"`
	ExpiresAt   string   `json:"expires_at,omitempty"`
}

// CreateAPIKeyDTO is a swagger-friendly response for the create API key endpoint.
type CreateAPIKeyDTO struct {
	ApiKey *APIKeyDTO `json:"api_key,omitempty"`
	Secret string     `json:"secret,omitempty"`
}

// ListAPIKeysDTO is a swagger-friendly response for the list API keys endpoint.
type ListAPIKeysDTO struct {
	ApiKeys []*APIKeyDTO `json:"api_keys,omitempty"`
	Total   int64        `json:"total,omitempty"`
	Page    int32        `json:"page,omitempty"`
	Limit   int32        `json:"limit,omitempty"`
}

// GetAPIKeyResponseSwagger is a swagger-friendly response for get API key.
type GetAPIKeyDTO struct {
	ApiKey *APIKeyDTO `json:"api_key,omitempty"`
}

// RevokeAPIKeyResponseSwagger is a swagger-friendly response for the revoke API key endpoint.
type RevokeAPIKeyDTO struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
}

// RotateAPIKeyResponseSwagger is a swagger-friendly response for the rotate API key endpoint.
type RotateAPIKeyDTO struct {
	ApiKey *APIKeyDTO `json:"api_key,omitempty"`
	Secret string     `json:"secret,omitempty"`
}

// ValidateAPIKeyRequestSwagger is a swagger-friendly request body for validating an API key.
type ValidateAPIKeyRequest struct {
	Key string `json:"key,omitempty"`
}

// ValidateAPIKeyResponseSwagger is a swagger-friendly response for the validate API key endpoint.
type ValidateAPIKeyDTO struct {
	Valid  bool       `json:"valid,omitempty"`
	ApiKey *APIKeyDTO `json:"api_key,omitempty"`
}
