package conv

import (
	"database/sql"
	"time"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infrastructures/config"
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
		EmailVerifiedAt: false,
		CreatedAt:       time.Now().Unix(),
		CreatedBy:       id,
		UpdatedAt:       time.Now().Unix(),
		UpdatedBy:       sql.NullString{},
		DeletedAt:       sql.NullInt64{},
		DeletedBy:       sql.NullString{},
	}

	endPoint := helper.EndPointMarshal()

	access := &model.Access{
		AppId:          req.AppId,
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

func RegisterModelToResp(user *model.User, emailFormat string) *dto.UserResp {
	phoneNumber := ""
	if user.PhoneNumber.Valid {
		phoneNumber = user.PhoneNumber.String
	} else {
		phoneNumber = "null"
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
		EmailVerifiedAt: false,
	}

	return userResp
}
