package usecase

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infrastructures/config"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/helper/conv"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/encryption"
	_error "github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/error"
)

type Oauth2UsecaseImpl struct {
	userRepo   repository.UserSqlRepo
	oauth2Repo repository.Oauth2ProviderRepo
}

func NewOauth2UsecaseImpl(
	userRepo repository.UserSqlRepo,
	oauth2Repo repository.Oauth2ProviderRepo,
) usecase.Oauth2Usecase {
	return &Oauth2UsecaseImpl{
		userRepo:   userRepo,
		oauth2Repo: oauth2Repo,
	}
}

func (o *Oauth2UsecaseImpl) GoogleClaimUser(ctx context.Context, req *dto.LoginGoogleReq) (*dto.LoginGoogleResp, error) {
	var googleOauthToken *model.GoogleOauth2Token

	if req.Device == "mobile" {
		googleToken, err := encryption.DecryptStringCBC(req.Token, config.AesCBC, config.AesCBCIV)
		if err != nil {
			log.Err(err).Msg("invalid decrypt token")
			return nil, _error.ErrStringDefault(http.StatusForbidden)
		}

		err = json.Unmarshal([]byte(googleToken), &googleOauthToken)
		if err != nil {

			return nil, _error.ErrStringDefault(http.StatusForbidden)
		}
	} else {
		googleCode, err := encryption.DecryptStringCBC(req.Token, config.AesCBC, config.AesCBCIV)
		if err != nil {
			return nil, _error.ErrStringDefault(http.StatusForbidden)
		}

		googleOauthToken, err = o.oauth2Repo.GetGoogleOauthToken(googleCode)
		if err != nil {
			return nil, _error.ErrStringDefault(http.StatusForbidden)
		}
	}

	googleUser, err := o.oauth2Repo.GetGoogleOauthUser(googleOauthToken)
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusForbidden)
	}

	err = o.userRepo.OpenConn(ctx)
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer o.userRepo.CloseConn()

	exist, err := o.userRepo.CheckUserByEmail(ctx, googleUser.Email)
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	googleUserResp := conv.GoogleClaimModelToResp(googleUser, exist)
	return googleUserResp, nil
}
