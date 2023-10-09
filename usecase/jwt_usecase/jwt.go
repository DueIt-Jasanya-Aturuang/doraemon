package jwt_usecase

import (
	"time"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra"
)

type Jwt struct {
	UserID string
	Key    string
	Exp    time.Duration
}

func (j *Jwt) AccessTokenDefault(userID string) *Jwt {
	return &Jwt{
		UserID: userID,
		Key:    infra.AccessTokenKeyHS,
		Exp:    infra.AccessTokenKeyExpHS,
	}
}

func (j *Jwt) RefreshTokenDefault(userID string, rememberMe bool) *Jwt {
	var exp time.Duration
	if rememberMe {
		exp = infra.RememberMeTokenExp
	} else {
		exp = infra.RefreshTokenKeyExpHS
	}

	return &Jwt{
		UserID: userID,
		Key:    infra.RefreshTokenKeyHS,
		Exp:    exp,
	}
}

func (j *Jwt) ForgotPasswordTokenDefault(userID string) *Jwt {
	return &Jwt{
		UserID: userID,
		Key:    infra.DefaultKey,
		Exp:    infra.ForgotPasswordTokenExp,
	}
}
