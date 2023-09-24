package domain

import (
	"context"
	"database/sql"
)

//counterfeiter:generate -o ./../mocks . UserSqlRepo
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

//counterfeiter:generate -o ./../mocks . UserUsecase
type UserUsecase interface {
	ChangePassword(ctx context.Context, req *RequestChangePassword) error
	ChangeUsername(ctx context.Context, req *RequestChangeUsername) error
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
	ID              string  `json:"id"`
	FullName        string  `json:"full_name"`
	Gender          string  `json:"gender"`
	Image           string  `json:"image"`
	Username        string  `json:"username"`
	Email           string  `json:"email"`
	EmailFormat     string  `json:"email_format"`
	PhoneNumber     *string `json:"phone_number"`
	EmailVerifiedAt bool    `json:"activited"`
}

type ResponseActivasiAccount struct {
	EmailVerifiedAt bool `json:"activited"`
}

type RequestChangePassword struct {
	OldPassword string `json:"old_password"`
	Password    string `json:"password"`
	RePassword  string `json:"re_password"`
	UserID      string // UserID get in header
}

type RequestChangeUsername struct {
	Username string `json:"username"`
	UserID   string // UserID get in header
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

var CheckUserByEmail string
var CheckUserByUsername string
var GetUserByID string
var GetUserByEmail string
var GetUserByUsername string

// func NewUserRepositoryGet(urgss UserRepositoryGetSqlSelect) string {
// 	switch urgss {
// 	case GetUserByID:
// 		return `SELECT id, fullname, gender, image, username, email, password, phone_number, email_verified_at,
//        				 created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
// 			  FROM m_users WHERE id = $1`
// 	case GetUserByEmail:
// 		return `SELECT id, fullname, gender, image, username, email, password, phone_number, email_verified_at,
//        				 created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
// 			  FROM m_users WHERE email = $1`
// 	case GetUserByUsername:
// 		return `SELECT id, fullname, gender, image, username, email, password, phone_number, email_verified_at,
//        				 created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
// 			  FROM m_users WHERE username = $1`
// 	case GetUserByEmailOrUsername:
// 		return `SELECT id, fullname, gender, image, username, email, password, phone_number, email_verified_at,
//        				 created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
// 			  FROM m_users WHERE username = $1 OR email = $2`
// 	}
//
// 	return `SELECT id, fullname, gender, image, username, email, password, phone_number, email_verified_at,
//        				 created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
// 			  FROM m_users WHERE id = $1`
// }

// func NewUserRepositoryCheck(urgss UserRepositoryGetSqlSelect) string {
// 	switch urgss {
// 	case CheckActivasiUserByID:
// 		return "SELECT email_verified_at FROM m_users WHERE id = $1"
// 	case CheckUserByEmail:
// 		return "SELECT EXISTS(SELECT 1 FROM m_users WHERE email = $1 AND deleted_at IS NULL)"
// 	case CheckUserByUsername:
// 		return "SELECT EXISTS(SELECT 1 FROM m_users WHERE username = $1 AND deleted_at IS NULL)"
// 	}
//
// 	return ""
// }
