package usecase

import (
	"context"
	"time"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/usecase"
)

type AuthUsecaseImpl struct {
	userRepo   repository.UserSqlRepo
	oauth2Repo repository.Oauth2ProviderRepo
	appRepo    repository.AppSqlRepo
	accessRepo repository.AccessSqlRepo
	timeout    time.Duration
}

func NewAuthUsecaseImpl(
	userRepo repository.UserSqlRepo,
	oauth2Repo repository.Oauth2ProviderRepo,
	appRepo repository.AppSqlRepo,
	accessRepo repository.AccessSqlRepo,
	timeout time.Duration,
) usecase.AuthUsecase {
	return &AuthUsecaseImpl{
		userRepo:   userRepo,
		oauth2Repo: oauth2Repo,
		appRepo:    appRepo,
		accessRepo: accessRepo,
		timeout:    timeout,
	}
}

func (a *AuthUsecaseImpl) Login(ctx context.Context, req *dto.LoginReq) (*dto.UserResp, error) {
	// TODO implement me
	panic("implement me")
}

func (a *AuthUsecaseImpl) Logout(ctx context.Context, req *dto.LogoutReq) error {
	// TODO implement me
	panic("implement me")
}

func (a *AuthUsecaseImpl) Register(ctx context.Context, req *dto.RegisterReq) (*dto.UserResp, error) {
	// TODO implement me
	panic("implement me")
}
