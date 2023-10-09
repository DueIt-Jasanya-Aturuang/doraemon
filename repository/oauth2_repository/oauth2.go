package oauth2_repository

import (
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository"
)

type Oauth2RepositoryImpl struct {
	clientID    string
	secretID    string
	redirectURI string
}

func NewOauth2RepositoryImpl() repository.Oauth2Repository {
	return &Oauth2RepositoryImpl{
		clientID:    infra.OauthClientId,
		secretID:    infra.OauthClientSecret,
		redirectURI: infra.OauthClientRedirectURI,
	}
}
