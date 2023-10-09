package usecase

import (
	"context"
	"time"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository"
)

type UserUsecase interface {
	GetUserByID(ctx context.Context, id string) (*ResponseUser, error)
	ChangePassword(ctx context.Context, req *RequestChangePassword) error
	ChangeUsername(ctx context.Context, req *RequestChangeUsername) error
	ForgottenPassword(ctx context.Context, req *RequestForgottenPassword) (string, error)
	ResetForgottenPassword(ctx context.Context, req *RequestResetForgottenPassword) error
	ActivasiAccount(c context.Context, email string) (bool, error)
}

type RequestChangePassword struct {
	OldPassword string
	Password    string
	RePassword  string
	UserID      string
}

type RequestChangeUsername struct {
	Username string
	UserID   string
}

type RequestForgottenPassword struct {
	Email string
}

type RequestResetForgottenPassword struct {
	Email      string
	Token      string
	Password   string
	RePassword string
}

type ResponseUser struct {
	ID              string
	FullName        string
	Gender          string
	Image           string
	Username        string
	Email           string
	EmailFormat     string
	PhoneNumber     *string
	EmailVerifiedAt bool
}

func ChangePasswordRequestToModel(password string, userID string) *repository.User {
	return &repository.User{
		ID:       userID,
		Password: password,
		AuditInfo: repository.AuditInfo{
			UpdatedAt: time.Now().Unix(),
			UpdatedBy: repository.NewNullString(userID),
		},
	}
}

func ChangeUsernameRequestToModel(username string, userID string) *repository.User {
	return &repository.User{
		ID:       userID,
		Username: username,
		AuditInfo: repository.AuditInfo{
			UpdatedAt: time.Now().Unix(),
			UpdatedBy: repository.NewNullString(userID),
		},
	}
}

func ActivasiAccountRequestToModel(userID string) *repository.User {
	user := &repository.User{
		ID:              userID,
		EmailVerifiedAt: true,
		AuditInfo: repository.AuditInfo{
			UpdatedAt: time.Now().Unix(),
			UpdatedBy: repository.NewNullString(userID),
		},
	}

	return user
}

func UserModelToResponse(u *repository.User) *ResponseUser {
	return &ResponseUser{
		ID:              u.ID,
		FullName:        u.FullName,
		Gender:          u.Gender,
		Image:           u.Image,
		Username:        u.Username,
		Email:           u.Email,
		EmailFormat:     FormatEmail(u.Email),
		PhoneNumber:     repository.GetNullString(u.PhoneNumber),
		EmailVerifiedAt: u.EmailVerifiedAt,
	}
}
