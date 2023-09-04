package usecase

import (
	"context"
	"encoding/json"
	"time"

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
	userRepo     repository.UserSqlRepo
	oauth2Repo   repository.Oauth2ProviderRepo
	appRepo      repository.AppSqlRepo
	accessRepo   repository.AccessSqlRepo
	securityRepo repository.SecuritySqlRepo
	timeout      time.Duration
}

func NewOauth2UsecaseImpl(
	userRepo repository.UserSqlRepo,
	oauth2Repo repository.Oauth2ProviderRepo,
	appRepo repository.AppSqlRepo,
	accessRepo repository.AccessSqlRepo,
	securityRepo repository.SecuritySqlRepo,
	timeout time.Duration,
) usecase.Oauth2Usecase {
	return &Oauth2UsecaseImpl{
		userRepo:     userRepo,
		oauth2Repo:   oauth2Repo,
		appRepo:      appRepo,
		accessRepo:   accessRepo,
		securityRepo: securityRepo,
	}
}

func (o *Oauth2UsecaseImpl) GoogleClaimUser(c context.Context, req *dto.LoginGoogleReq) (*dto.LoginGoogleResp, error) {
	ctx, cancel := context.WithTimeout(c, o.timeout)
	defer cancel()

	var googleOauthToken *model.GoogleOauth2Token

	if req.Device == "mobile" {
		googleToken, err := encryption.DecryptStringCBC(req.Token, config.AesCBC, config.AesCBCIV)
		if err != nil {
			return nil, _error.ErrString("UNATHORIZATION", 401)
		}

		err = json.Unmarshal([]byte(googleToken), &googleOauthToken)
		if err != nil {
			return nil, _error.ErrString("UNATHORIZATION", 401)
		}
	} else {
		googleCode, err := encryption.DecryptStringCBC(req.Token, config.AesCBC, config.AesCBCIV)
		if err != nil {
			return nil, _error.ErrString("UNATHORIZATION", 401)
		}

		googleOauthToken, err = o.oauth2Repo.GetGoogleOauthToken(googleCode)
		if err != nil {
			return nil, _error.ErrString("UNATHORIZATION", 401)
		}
	}

	googleUser, err := o.oauth2Repo.GetGoogleOauthUser(googleOauthToken)
	if err != nil {
		return nil, _error.ErrString("UNATHORIZATION", 401)
	}

	exist, err := o.userRepo.CheckUserByEmail(ctx, googleUser.Email)
	if err != nil {
		return nil, _error.ErrString("INTERNAL SERVER ERROR", 500)
	}

	googleUserResp := conv.GoogleClaimModelToResp(googleUser, exist)
	return googleUserResp, nil
}
