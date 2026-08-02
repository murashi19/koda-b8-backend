package users

import (
	"time"

	"github.com/murashi19/koda-b8-backend/internal/models"
)

type UpdateUserRequest struct {
	Email       string         `json:"email" binding:"required,email"`
	FullName    string         `json:"full_name" binding:"required"`
	PhoneNumber *string        `json:"phone_number"`
	Avatar      *string        `json:"avatar"`
	BirthDate   *time.Time     `json:"birth_date"`
	Gender      *models.Gender `json:"gender"`
}
