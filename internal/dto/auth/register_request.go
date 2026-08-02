package auth

type RegisterRequest struct {
	FullName string `json:"full_name" binding:"required,min=3,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}
