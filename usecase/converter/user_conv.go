package converter

import (
	"time"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase/helper"
)

func ChangePasswordReqToModel(password string, userID string) *domain.User {
	return &domain.User{
		ID:       userID,
		Password: password,
		AuditInfo: domain.AuditInfo{
			UpdatedAt: time.Now().Unix(),
			UpdatedBy: helper.NewNullString(userID),
		},
	}
}

func ChangeUsernameReqToModel(username string, userID string) *domain.User {
	return &domain.User{
		ID:       userID,
		Username: username,
		AuditInfo: domain.AuditInfo{
			UpdatedAt: time.Now().Unix(),
			UpdatedBy: helper.NewNullString(userID),
		},
	}
}

func ActivasiAccountReqToModel(userID string) *domain.User {
	user := &domain.User{
		ID:              userID,
		EmailVerifiedAt: true,
		AuditInfo: domain.AuditInfo{
			UpdatedAt: time.Now().Unix(),
			UpdatedBy: helper.NewNullString(userID),
		},
	}

	return user
}

func UserModelToResp(u *domain.User) *domain.ResponseUser {
	return &domain.ResponseUser{
		ID:              u.ID,
		FullName:        u.FullName,
		Gender:          u.Gender,
		Image:           u.Image,
		Username:        u.Username,
		Email:           u.Email,
		EmailFormat:     helper.EmailFormat(u.Email),
		PhoneNumber:     helper.GetNullString(u.PhoneNumber),
		EmailVerifiedAt: u.EmailVerifiedAt,
	}
}
