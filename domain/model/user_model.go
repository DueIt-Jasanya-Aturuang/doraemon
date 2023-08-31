package model

import "database/sql"

type User struct {
	ID              string         `json:"user_id"`
	FullName        string         `json:"full_name"`
	Gender          string         `json:"gender"`
	Image           string         `json:"string"`
	Username        string         `json:"username"`
	Email           string         `json:"email"`
	Password        string         `json:"password"`
	PhoneNumber     sql.NullString `json:"phone_number"`
	EmailVerifiedAt bool           `json:"email_verified_at"`
	CreatedAt       int64          `json:"created_at"`
	CreatedBy       string         `json:"created_by"`
	UpdatedAt       int64          `json:"updated_at"`
	UpdatedBy       sql.NullString `json:"updated_by"`
	DeletedAt       sql.NullInt64  `json:"deleted_at"`
	DeletedBy       sql.NullString `json:"deleted_by"`
}
