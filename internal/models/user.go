package models

import "time"

type UserRole string

const (
	RoleAdmin    UserRole = "ADMIN"
	RoleCustomer UserRole = "CUSTOMER"
)

type User struct {
	ID         int64     `db:"id"`
	Email      string    `db:"email"`
	Password   string    `db:"password"`
	Role       UserRole  `db:"role"`
	Fullname   string    `db:"full_name"`
	IsVerified bool      `db:"is_verified"`
	IsActive   bool      `db:"is_active"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
type UserDetail struct {
	User
	Profile *UserProfile
}
