package repository

import (
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
)

type GoogleOauthRepo interface {
	GetGoogleOauthToken(code string) (*model.GoogleOauthToken, error)
	GetGoogleOauthUser(token *model.GoogleOauthToken) (*model.GoogleOauthUser, error)
}
