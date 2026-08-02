package users

import (
	"time"

	"github.com/murashi19/koda-b8-backend/internal/models"
)

type UpdateProfileRequest struct {
	FullName    string         `json:"full_name" binding:"required,min=3,max=100"`
	Email       string         `json:"email" binding:"required,email"`
	PhoneNumber *string        `json:"phone_number"`
	BirthDate   *time.Time     `json:"birth_date"`
	Gender      *models.Gender `json:"gender"`
}
