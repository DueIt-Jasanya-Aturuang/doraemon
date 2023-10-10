package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

type AuthUsecase interface {
	Login(ctx context.Context, req *RequestLogin) (*ResponseAuth, error)
	Register(ctx context.Context, req *RequestRegister) (*ResponseAuth, error)
}

type ResponseAuth struct {
	User    *ResponseUser
	Profile *ResponseProfileDueit
	Token   *ResponseJWT
}

type RequestRegister struct {
	FullName        string
	Username        string
	Email           string
	Password        string
	RePassword      string
	EmailVerifiedAt bool
	AppID           string
	Role            int8
	RememberMe      bool
}

type RequestLogin struct {
	EmailOrUsername string
	Password        string
	RememberMe      bool
	Oauth2          bool
	AppID           string
}

type AuthMergeModelToResponse struct {
	User        *repository.User
	Profile     *ResponseProfileDueit
	Token       string
	FormatEmail string
}

func (m *AuthMergeModelToResponse) Execute() *ResponseAuth {
	user := &ResponseUser{
		ID:              m.User.ID,
		FullName:        m.User.FullName,
		Gender:          m.User.Gender,
		Image:           m.User.Image,
		Username:        m.User.Username,
		Email:           m.User.Email,
		EmailFormat:     m.FormatEmail,
		PhoneNumber:     repository.GetNullString(m.User.PhoneNumber),
		EmailVerifiedAt: m.User.EmailVerifiedAt,
	}

	profile := &ResponseProfileDueit{
		ProfileID: m.Profile.ProfileID,
		Quote:     m.Profile.Quote,
		Profesi:   m.Profile.Profesi,
	}

	token := &ResponseJWT{
		Token: m.Token,
	}

	resp := &ResponseAuth{
		User:    user,
		Profile: profile,
		Token:   token,
	}
	return resp
}

func AuthRegisterRequestToModel(req *RequestRegister) (*repository.User, *repository.Access) {
	endpoint, err := json.Marshal(util.Endpoint)
	if err != nil {
		log.Warn().Msgf(util.LogErrMarshal, util.Endpoint, err)
	}

	id := util.NewUUID
	user := &repository.User{
		ID:              id,
		FullName:        req.FullName,
		Gender:          "undefined",
		Image:           infra.DefaultImage,
		Username:        req.Username,
		Email:           req.Email,
		Password:        req.Password,
		PhoneNumber:     sql.NullString{},
		EmailVerifiedAt: req.EmailVerifiedAt,
		AuditInfo: repository.AuditInfo{
			CreatedAt: time.Now().Unix(),
			CreatedBy: id,
			UpdatedAt: time.Now().Unix(),
		},
	}

	access := &repository.Access{
		AppId:          req.AppID,
		UserId:         id,
		RoleId:         req.Role,
		AccessEndpoint: string(endpoint),
		AuditInfo: repository.AuditInfo{
			CreatedAt: time.Now().Unix(),
			CreatedBy: id,
			UpdatedAt: time.Now().Unix(),
		},
	}
	return user, access
}
