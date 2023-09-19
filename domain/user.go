package domain

import (
	"context"
	"database/sql"
)

//counterfeiter:generate -o ./../mocks . UserSqlRepo
type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	UpdateActivasiUser(ctx context.Context, user *User) error
	UpdatePasswordUser(ctx context.Context, user *User) error
	CheckActivasiUserByID(ctx context.Context, id string) (bool, error)
	CheckUserByEmail(ctx context.Context, email string) (bool, error)
	CheckUserByUsername(ctx context.Context, username string) (bool, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByEmailOrUsername(ctx context.Context, emailOrUsername string) (*User, error)
	UnitOfWorkRepository
}

//counterfeiter:generate -o ./../mocks . UserUsecase
type UserUsecase interface {
	ResetPassword(ctx context.Context, req *RequestResetPassword) error
	ForgottenPassword(ctx context.Context, req *RequestForgottenPassword) (string, error)
	ResetForgottenPassword(ctx context.Context, req *RequestResetForgottenPassword) error
	ActivasiAccount(c context.Context, email string) (resp *ResponseActivasiAccount, err error)
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

type ResponseUser struct {
	ID              string `json:"id"`
	FullName        string `json:"full_name"`
	Gender          string `json:"gender"`
	Image           string `json:"image"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	EmailFormat     string `json:"email_format"`
	PhoneNumber     string `json:"phone_number"`
	EmailVerifiedAt bool   `json:"activited"`
}

type ResponseActivasiAccount struct {
	EmailVerifiedAt bool `json:"activited"`
}

type RequestResetPassword struct {
	OldPassword string `json:"old_password"`
	Password    string `json:"password"`
	RePassword  string `json:"re_password"`
	UserID      string // UserID get in header
}

type RequestForgottenPassword struct {
	Email string // Email get in query param
}

type RequestResetForgottenPassword struct {
	Email      string // Email get in query param
	Token      string // Token get in query param
	Password   string `json:"password"`
	RePassword string `json:"re_password"`
}
