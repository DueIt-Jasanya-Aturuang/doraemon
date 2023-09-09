package model

import (
	"time"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/infrastructures/config"
)

type Jwt struct {
	UserID string
	Key    string
	Exp    time.Duration
}

func (j *Jwt) AccessTokenDefault(userID string) *Jwt {
	return &Jwt{
		UserID: userID,
		Key:    config.AccessTokenKeyHS,
		Exp:    config.AccessTokenKeyExpHS,
	}
}

func (j *Jwt) RefreshTokenDefault(userID string, rememberMe bool) *Jwt {
	var exp time.Duration
	if rememberMe {
		exp = config.RememberMeTokenExp
	} else {
		exp = config.RefreshTokenKeyExpHS
	}

	return &Jwt{
		UserID: userID,
		Key:    config.RefreshTokenKeyHS,
		Exp:    exp,
	}
}

func (j *Jwt) ForgotPasswordTokenDefault(userID string) *Jwt {
	return &Jwt{
		UserID: userID,
		Key:    config.DefaultKey,
		Exp:    config.ForgotPasswordTokenExp,
	}
}
