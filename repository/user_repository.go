package repository

import (
	"context"
	"database/sql"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	UpdateActivasi(ctx context.Context, user *User) error
	UpdatePassword(ctx context.Context, user *User) error
	UpdateUsername(ctx context.Context, user *User) error
	CheckActivasiUser(ctx context.Context, id string) (bool, error)
	GetByEmailOrUsername(ctx context.Context, s string) (*User, error)
	Check(ctx context.Context, s string) (bool, error)
	Get(ctx context.Context, s string) (*User, error)
	UnitOfWorkRepository
}

type User struct {
	ID              string
	FullName        string
	Gender          string
	Image           string
	Username        string
	Email           string
	Password        string
	PhoneNumber     sql.NullString
	EmailVerifiedAt bool
	AuditInfo
}

var CheckUserByEmail string
var CheckUserByUsername string
var GetUserByID string
var GetUserByEmail string
var GetUserByUsername string
