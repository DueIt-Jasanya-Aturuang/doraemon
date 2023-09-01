package model

import (
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infrastructures/config"
	"time"
)

type Jwt struct {
	UUID       string
	Type       string
	RememberMe bool
	Key        string
	Exp        time.Duration
}

func (j *Jwt) AccessTokenDefault(uuid string, rememberMe bool) *Jwt {
	return &Jwt{
		UUID:       uuid,
		Type:       "access_token",
		RememberMe: rememberMe,
		Key:        config.AccessTokenKeyHS,
		Exp:        config.AccessTokenKeyExpHS,
	}
}

func (j *Jwt) RefreshTokenDefault(uuid string, rememberMe bool) *Jwt {
	var exp time.Duration
	if rememberMe {
		exp = config.RememberMeTokenExp
	} else {
		exp = config.RefreshTokenKeyExpHS
	}

	return &Jwt{
		UUID:       uuid,
		Type:       "refresh_token",
		RememberMe: rememberMe,
		Key:        config.RefreshTokenKeyHS,
		Exp:        exp,
	}
}

func (j *Jwt) ForgotPasswordTokenDefault(uuid string) *Jwt {
	return &Jwt{
		UUID:       uuid,
		Type:       "forgot_password",
		RememberMe: false,
		Key:        config.DefaultKey,
		Exp:        config.ForgotPasswordTokenExp,
	}
}
