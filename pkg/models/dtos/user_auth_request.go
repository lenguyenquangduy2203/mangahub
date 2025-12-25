package dtos

// UserCreationRequest represents a standard user creation request
type UserAuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}
