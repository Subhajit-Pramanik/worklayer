package apikey

import (
	"time"

	pb "github.com/vyolayer/vyolayer/proto/apikey/v1"
)

func protoApiKeyToDTO(apiKey *pb.APIKey) *APIKeyDTO {
	apiKeyDTO := &APIKeyDTO{
		Id:             apiKey.GetId(),
		OrganizationId: apiKey.GetOrganizationId(),
		ProjectId:      apiKey.GetProjectId(),
		Name:           apiKey.GetName(),
		Description:    apiKey.GetDescription(),
		Prefix:         apiKey.GetPrefix(),
		Environment:    apiKey.GetEnvironment(),
		Status:         apiKey.GetStatus(),
		Scopes:         apiKey.GetScopes(),
		CreatedBy:      apiKey.GetCreatedBy(),
		LastUsedIp:     apiKey.GetLastUsedIp(),
		LastUsedUa:     apiKey.GetLastUsedUa(),
		RevokedBy:      apiKey.GetRevokedBy(),
		CreatedAt:      apiKey.GetCreatedAt().AsTime().Format(time.RFC3339),
		UpdatedAt:      apiKey.GetUpdatedAt().AsTime().Format(time.RFC3339),
	}

	if apiKey.GetLastUsedAt() != nil {
		apiKeyDTO.LastUsedAt = apiKey.GetLastUsedAt().AsTime().Format(time.RFC3339)
	}

	if apiKey.GetExpiresAt() != nil {
		apiKeyDTO.ExpiresAt = apiKey.GetExpiresAt().AsTime().Format(time.RFC3339)
	}

	if apiKey.GetRevokedAt() != nil {
		apiKeyDTO.RevokedAt = apiKey.GetRevokedAt().AsTime().Format(time.RFC3339)
	}

	if len(apiKey.GetScopes()) == 0 {
		apiKeyDTO.Scopes = []string{}
	}

	return apiKeyDTO
}
