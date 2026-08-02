package auth

import (
	"github.com/murashi19/koda-b8-backend/internal/dto/users"
)

type LoginResponse struct {
	Token string             `json:"token"`
	User  users.UserResponse `json:"user"`
}
