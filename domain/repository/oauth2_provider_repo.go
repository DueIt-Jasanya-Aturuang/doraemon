package repository

import (
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
)

//counterfeiter:generate -o ./../mocks . Oauth2ProviderRepo
type Oauth2ProviderRepo interface {
	GetGoogleOauthToken(code string) (*model.GoogleOauth2Token, error)
	GetGoogleOauthUser(token *model.GoogleOauth2Token) (*model.GoogleOauth2User, error)
}
