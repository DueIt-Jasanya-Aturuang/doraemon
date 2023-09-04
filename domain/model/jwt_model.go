package model

import (
	"time"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/infrastructures/config"
)

type Jwt struct {
	TokenID    string
	UserID     string
	Type       string
	RememberMe bool
	Key        string
	Exp        time.Duration
}

func (j *Jwt) AccessTokenDefault(tokenID string, userID string, rememberMe bool) *Jwt {
	return &Jwt{
		TokenID:    tokenID,
		UserID:     userID,
		Type:       "access_token",
		RememberMe: rememberMe,
		Key:        config.AccessTokenKeyHS,
		Exp:        config.AccessTokenKeyExpHS,
	}
}

func (j *Jwt) RefreshTokenDefault(tokenID string, userID string, rememberMe bool) *Jwt {
	var exp time.Duration
	if rememberMe {
		exp = config.RememberMeTokenExp
	} else {
		exp = config.RefreshTokenKeyExpHS
	}

	return &Jwt{
		TokenID:    tokenID,
		UserID:     userID,
		Type:       "refresh_token",
		RememberMe: rememberMe,
		Key:        config.RefreshTokenKeyHS,
		Exp:        exp,
	}
}

func (j *Jwt) ForgotPasswordTokenDefault(tokenID string, userID string) *Jwt {
	return &Jwt{
		TokenID:    tokenID,
		UserID:     userID,
		Type:       "forgot_password",
		RememberMe: false,
		Key:        config.DefaultKey,
		Exp:        config.ForgotPasswordTokenExp,
	}
}
