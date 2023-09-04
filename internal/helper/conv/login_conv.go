package conv

import (
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
)

func LoginModelToResp(user *model.User, accessToken string, emailFormat string) *dto.UserResp {
	phoneNumber := ""
	if user.PhoneNumber.Valid {
		phoneNumber = user.PhoneNumber.String
	} else {
		phoneNumber = "null"
	}

	return &dto.UserResp{
		ID:              user.ID,
		FullName:        user.FullName,
		Gender:          user.Gender,
		Image:           user.Image,
		Username:        user.Username,
		Email:           user.Email,
		EmailFormat:     emailFormat,
		PhoneNumber:     phoneNumber,
		EmailVerifiedAt: user.EmailVerifiedAt,
		Token:           accessToken,
	}
}
