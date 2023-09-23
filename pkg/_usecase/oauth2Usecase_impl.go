package _usecase

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/pkg/converter"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util/encryption"
)

type Oauth2UsecaseImpl struct {
	userRepo   domain.UserRepository
	oauth2Repo domain.Oauth2Repository
}

func NewOauth2UsecaseImpl(
	userRepo domain.UserRepository,
	oauth2Repo domain.Oauth2Repository,
) domain.Oauth2Usecase {
	return &Oauth2UsecaseImpl{
		userRepo:   userRepo,
		oauth2Repo: oauth2Repo,
	}
}

func (o *Oauth2UsecaseImpl) GoogleClaimUser(ctx context.Context, req *domain.RequestLoginWithGoogle) (*domain.ResponseLoginWithGoogle, error) {
	var googleOauthToken *domain.Oauth2GoogleToken

	if req.Device == "mobile" {
		googleToken, err := encryption.DecryptStringCBC(req.Token, infra.AesCBC, infra.AesCBCIV)
		if err != nil {
			log.Warn().Msgf("failed decrypt token | err : %v", err)
			return nil, InvalidTokenOauth
		}

		err = json.Unmarshal([]byte(googleToken), &googleOauthToken)
		if err != nil {
			log.Warn().Msgf("invalid request token encrypt | data : %v | source : %v | err : %v", googleToken, googleOauthToken, err)
			return nil, InvalidTokenOauth
		}
	} else {
		googleCode, err := encryption.DecryptStringCBC(req.Token, infra.AesCBC, infra.AesCBCIV)
		if err != nil {
			log.Warn().Msgf("failed decrypt token | err : %v", err)
			return nil, InvalidTokenOauth
		}

		googleOauthToken, err = o.oauth2Repo.GetGoogleToken(googleCode)
		if err != nil {
			return nil, InvalidTokenOauth
		}
	}

	googleUser, err := o.oauth2Repo.GetGoogleUser(googleOauthToken)
	if err != nil {
		return nil, err
	}

	err = o.userRepo.OpenConn(ctx)
	if err != nil {
		return nil, err
	}
	defer o.userRepo.CloseConn()

	domain.CheckUserByEmail = googleUser.Email
	exist, err := o.userRepo.Check(ctx, domain.CheckUserByEmail)
	if err != nil {
		return nil, err
	}

	googleUserResp := converter.GoogleClaimModelToResp(googleUser, exist)
	return googleUserResp, nil
}
