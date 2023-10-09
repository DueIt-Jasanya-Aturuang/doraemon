package oauth2_usecase

import (
	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase"
)

type Oauth2UsecaseImpl struct {
	userRepo   repository.UserRepository
	oauth2Repo repository.Oauth2Repository
}

func NewOauth2UsecaseImpl(
	userRepo repository.UserRepository,
	oauth2Repo repository.Oauth2Repository,
) usecase.Oauth2Usecase {
	return &Oauth2UsecaseImpl{
		userRepo:   userRepo,
		oauth2Repo: oauth2Repo,
	}
}
