package dto

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name" validate:"required,min=3,max=100"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required"`
}

type UpdateAccountRequest struct {
	FullName string `json:"full_name,omitempty" validate:"omitempty,min=3,max=100"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
}

type AccountResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	FullName  string `json:"full_name"`
	Role      string `json:"role"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type LoginResponse struct {
	Account      AccountResponse `json:"account"`
	AccessToken  string          `json:"access_token"`
	RefreshToken string          `json:"refresh_token,omitempty"`
	TokenType    string          `json:"token_type"`
	ExpiresIn    int             `json:"expires_in"`
}

type ListAccountResponse struct {
	Data       []AccountResponse `json:"data"`
	Total      int64             `json:"total"`
	Limit      int               `json:"limit"`
	Offset     int               `json:"offset"`
	TotalPages int               `json:"total_pages"`
}
