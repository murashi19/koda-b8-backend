package models

import "time"

type RefreshToken struct {
	ID        int64      `db:"id" json:"id"`
	UserID    int64      `db:"user_id" json:"user_id"`
	TokenHash string     `db:"token_hash" json:"-"`
	ExpiresAt time.Time  `db:"expires_at" json:"expires_at"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	RevokedAt *time.Time `db:"revoked_at" json:"-"`
	UserAgent *string    `db:"user_agent" json:"-"`
	IPAddress *string    `db:"ip_address" json:"-"`
}
