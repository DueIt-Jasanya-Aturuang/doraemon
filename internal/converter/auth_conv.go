package converter

import (
	"database/sql"
	"time"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/helper"
)

func AuthModelToResp(user *domain.User, profile *domain.Profile, emailFormat string) *domain.ResponseAuth {
	resp := &domain.ResponseAuth{
		ResponseUser: domain.ResponseUser{
			ID:              user.ID,
			FullName:        user.FullName,
			Gender:          user.Gender,
			Image:           user.Image,
			Username:        user.Username,
			Email:           user.Email,
			EmailFormat:     emailFormat,
			PhoneNumber:     helper.GetNullString(user.PhoneNumber),
			EmailVerifiedAt: user.EmailVerifiedAt,
		},
		ResponseProfile: domain.ResponseProfile{
			ProfileID: profile.ProfileID,
			Quote:     profile.Quote,
			Profesi:   profile.Profesi,
		},
		ResponseJwtToken: domain.ResponseJwtToken{},
	}

	return resp
}

func RegisterReqToModel(req *domain.RequestRegister, id string) (*domain.User, *domain.Access) {
	user := &domain.User{
		ID:              id,
		FullName:        req.FullName,
		Gender:          "undefined",
		Image:           infra.DefaultImage,
		Username:        req.Username,
		Email:           req.Email,
		Password:        req.Password,
		PhoneNumber:     sql.NullString{},
		EmailVerifiedAt: req.EmailVerifiedAt,
		AuditInfo: domain.AuditInfo{
			CreatedAt: time.Now().Unix(),
			CreatedBy: id,
			UpdatedAt: time.Now().Unix(),
		},
	}

	endPoint := helper.EndPointMarshal()

	access := &domain.Access{
		AppId:          req.AppID,
		UserId:         id,
		RoleId:         req.Role,
		AccessEndpoint: endPoint,
		AuditInfo: domain.AuditInfo{
			CreatedAt: time.Now().Unix(),
			CreatedBy: id,
			UpdatedAt: time.Now().Unix(),
		},
	}
	return user, access
}
