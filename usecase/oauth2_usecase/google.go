package oauth2_usecase

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util/encryption"
)

func (o *Oauth2UsecaseImpl) GoogleClaimUser(ctx context.Context, req *usecase.RequestLoginWithGoogle) (*usecase.ResponseLoginWithGoogle, error) {
	var googleOauthToken *repository.Oauth2GoogleToken

	if req.Device == "mobile" {
		googleToken, err := encryption.DecryptStringCBC(req.Token, infra.AesCBC, infra.AesCBCIV)
		if err != nil {
			log.Warn().Msgf("failed decrypt token | err : %v", err)
			return nil, usecase.InvalidTokenOauth
		}

		err = json.Unmarshal([]byte(googleToken), &googleOauthToken)
		if err != nil {
			log.Warn().Msgf("invalid request token encrypt | data : %v | source : %v | err : %v", googleToken, googleOauthToken, err)
			return nil, usecase.InvalidTokenOauth
		}
	} else {
		googleCode, err := encryption.DecryptStringCBC(req.Token, infra.AesCBC, infra.AesCBCIV)
		if err != nil {
			log.Warn().Msgf("failed decrypt token | err : %v", err)
			return nil, usecase.InvalidTokenOauth
		}

		googleOauthToken, err = o.oauth2Repo.GetGoogleToken(googleCode)
		if err != nil {
			return nil, usecase.InvalidTokenOauth
		}
	}

	googleUser, err := o.oauth2Repo.GetGoogleUser(googleOauthToken)
	if err != nil {
		return nil, err
	}

	repository.CheckUserByEmail = googleUser.Email
	exist, err := o.userRepo.Check(ctx, repository.CheckUserByEmail)
	if err != nil {
		return nil, err
	}

	googleUserResp := usecase.GoogleClaimModelToResp(googleUser, exist)
	return googleUserResp, nil
}
