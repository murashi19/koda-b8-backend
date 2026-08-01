package models

import "time"

type Gender string

const (
	GenderMale   Gender = "MALE"
	GenderFemale Gender = "FEMALE"
)

type UserProfile struct {
	UserID      int64      `db:"user_id"`
	FullName    string     `db:"full_name"`
	PhoneNumber *string    `db:"phone_number"`
	Avatar      *string    `db:"avatar"`
	BirthDate   *time.Time `db:"birth_date"`
	Gender      *Gender    `db:"gender"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}
