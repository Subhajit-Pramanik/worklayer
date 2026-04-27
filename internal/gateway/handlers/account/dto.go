package account

// RegisterRequestDTO is the swagger model for user registration.
type RegisterRequestDTO struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`
}

// UserDTO is the swagger model for a user.
type UserDTO struct {
	Id        string `json:"id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// RegisterResponseDTO is the swagger model for user registration response.
type RegisterResponseDTO struct {
	User    *UserDTO `json:"user,omitempty"`
	Message string   `json:"message,omitempty"`
}

// ResendVerificationEmailRequestDTO is the swagger model for resending verification email.
type ResendVerificationEmailRequestDTO struct {
	Email string `json:"email,omitempty"`
}

// LoginRequestDTO is the swagger model for user login.
type LoginRequestDTO struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

// LoginResponseDTO is the swagger model for login response.
type LoginResponseDTO struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int32  `json:"expires_in,omitempty"`
}

// SessionDTO is the swagger model for a session.
type SessionDTO struct {
	Id           string `json:"id,omitempty"`
	UserId       string `json:"user_id,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	UserAgent    string `json:"user_agent,omitempty"`
	ClientIp     string `json:"client_ip,omitempty"`
	IsBlocked    bool   `json:"is_blocked,omitempty"`
	ExpiresAt    string `json:"expires_at,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
}

// RefreshSessionResponseDTO is the swagger model for refresh session response.
type RefreshSessionResponseDTO struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int32  `json:"expires_in,omitempty"`
}

// AllSessionsResponseDTO is the swagger model for all sessions response.
type AllSessionsResponseDTO struct {
	Sessions []*SessionDTO `json:"sessions,omitempty"`
}

// ChangePasswordRequestDTO is the swagger model for changing password.
type ChangePasswordRequestDTO struct {
	CurrentPassword string `json:"current_password,omitempty"`
	NewPassword     string `json:"new_password,omitempty"`
}

// ForgotPasswordRequestDTO is the swagger model for forgot password.
type ForgotPasswordRequestDTO struct {
	Email string `json:"email,omitempty"`
}

// ResetPasswordRequestDTO is the swagger model for reset password.
type ResetPasswordRequestDTO struct {
	Password string `json:"password,omitempty"`
}
