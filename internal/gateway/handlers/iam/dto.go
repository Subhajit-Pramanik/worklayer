package iam

type RegisterRequest struct {
	Email    string `json:"email" example:"subhajit@vyolayer.com"`
	Password string `json:"password" example:"Password!123"`
	FullName string `json:"full_name" example:"Subhajit Pramanik"`
}

type ResendVerificationEmailRequest struct {
	Email string `json:"email" example:"subhajit@vyolayer.com"`
}

type LoginRequest struct {
	Email    string `json:"email" example:"subhajit@vyolayer.com"`
	Password string `json:"password" example:"Password!123"`
}

type AuthSessionDTO struct {
	AccessToken          string `json:"access_token"`
	AccessTokenExpiresAt int64  `json:"access_token_expires_at"`
}

type UserDTO struct {
	ID              string    `json:"id,omitempty"`
	Email           string    `json:"email,omitempty"`
	FullName        string    `json:"full_name,omitempty"`
	Status          string    `json:"status,omitempty"`
	IsEmailVerified bool      `json:"is_email_verified,omitempty"`
	JoinedAt        string    `json:"joined_at,omitempty"`
	Avatar          AvatarDTO `json:"avatar,omitzero"`
}

type AvatarDTO struct {
	ID            int64  `json:"id,omitempty"`
	Url           string `json:"url,omitempty"`
	FallbackChar  string `json:"fallback_char,omitempty"`
	FallbackColor string `json:"fallback_color,omitempty"`
}

type GetMeDTO struct {
	User *UserDTO `json:"user,omitempty"`
}
