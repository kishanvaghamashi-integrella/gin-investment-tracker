package dto

type CreateUserRequest struct {
	Name     string `json:"name,omitempty" binding:"required,min=3,max=50"`
	Email    string `json:"email,omitempty" binding:"required,email"`
	Password string `json:"password,omitempty" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email,omitempty" binding:"required,email"`
	Password string `json:"password,omitempty" binding:"required,min=6"`
}

type LoginResponse struct {
	ID    int64  `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Token string `json:"-"`
}

type GoogleUserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}
