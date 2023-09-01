package model

import (
	"time"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/infrastructures/config"
)

type Jwt struct {
	UUID       string
	UserID     string
	Type       string
	RememberMe bool
	Key        string
	Exp        time.Duration
}

func (j *Jwt) AccessTokenDefault(uuid string, userID string, rememberMe bool) *Jwt {
	return &Jwt{
		UUID:       uuid,
		UserID:     userID,
		Type:       "access_token",
		RememberMe: rememberMe,
		Key:        config.AccessTokenKeyHS,
		Exp:        config.AccessTokenKeyExpHS,
	}
}

func (j *Jwt) RefreshTokenDefault(uuid string, userID string, rememberMe bool) *Jwt {
	var exp time.Duration
	if rememberMe {
		exp = config.RememberMeTokenExp
	} else {
		exp = config.RefreshTokenKeyExpHS
	}

	return &Jwt{
		UUID:       uuid,
		UserID:     userID,
		Type:       "refresh_token",
		RememberMe: rememberMe,
		Key:        config.RefreshTokenKeyHS,
		Exp:        exp,
	}
}

func (j *Jwt) ForgotPasswordTokenDefault(uuid string, userID string) *Jwt {
	return &Jwt{
		UUID:       uuid,
		UserID:     userID,
		Type:       "forgot_password",
		RememberMe: false,
		Key:        config.DefaultKey,
		Exp:        config.ForgotPasswordTokenExp,
	}
}
