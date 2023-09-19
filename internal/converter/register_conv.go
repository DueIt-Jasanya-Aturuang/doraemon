package converter

import (
	"database/sql"
	"time"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra/config"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/helper"
)

func RegisterReqToModel(req *dto.RegisterReq, id string) (*model.User, *model.Access) {
	user := &model.User{
		ID:              id,
		FullName:        req.FullName,
		Gender:          "undefined",
		Image:           config.DefaultImage,
		Username:        req.Username,
		Email:           req.Email,
		Password:        req.Password,
		PhoneNumber:     sql.NullString{},
		EmailVerifiedAt: req.EmailVerifiedAt,
		CreatedAt:       time.Now().Unix(),
		CreatedBy:       id,
		UpdatedAt:       time.Now().Unix(),
		UpdatedBy:       sql.NullString{},
		DeletedAt:       sql.NullInt64{},
		DeletedBy:       sql.NullString{},
	}

	endPoint := helper.EndPointMarshal()

	access := &model.Access{
		AppId:          req.AppID,
		UserId:         id,
		RoleId:         req.Role,
		AccessEndpoint: endPoint,
		CreatedAt:      time.Now().Unix(),
		CreatedBy:      id,
		UpdatedAt:      time.Now().Unix(),
		UpdatedBy:      sql.NullString{},
		DeletedAt:      sql.NullInt64{},
		DeletedBy:      sql.NullString{},
	}
	return user, access
}

func RegisterModelToResp(user *model.User, profile *model.Profile, emailFormat string) (*dto.UserResp, *dto.ProfileResp) {
	phoneNumber := "null"
	if user.PhoneNumber.Valid {
		phoneNumber = user.PhoneNumber.String
	}

	userResp := &dto.UserResp{
		ID:              user.ID,
		FullName:        user.FullName,
		Gender:          user.Gender,
		Image:           user.Image,
		Username:        user.Username,
		Email:           user.Email,
		EmailFormat:     emailFormat,
		PhoneNumber:     phoneNumber,
		EmailVerifiedAt: user.EmailVerifiedAt,
	}

	profileResp := &dto.ProfileResp{
		ProfileID: profile.ProfileID,
		Quote:     profile.Quote,
		Profesi:   profile.Profesi,
	}
	return userResp, profileResp
}
