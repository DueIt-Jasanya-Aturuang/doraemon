package converter

import (
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
)

func LoginModelToResp(user *model.User, profile *model.Profile, emailFormat string) (*dto.UserResp, *dto.ProfileResp) {
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
