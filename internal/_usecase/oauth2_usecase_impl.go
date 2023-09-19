package _usecase

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/usecase"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra/config"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/converter"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util/encryption"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util/error"
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

	// checking condition apakah dia dari mobile atau web
	if req.Device == "mobile" {
		// kita decrypt request body dari user mobile make aes cbc
		googleToken, err := encryption.DecryptStringCBC(req.Token, config.AesCBC, config.AesCBCIV)
		if err != nil {
			log.Err(err).Msg("invalid decrypt token")
			return nil, _error.ErrStringDefault(http.StatusForbidden)
		}

		// dari hasil decrypt akan menghasilkan json yang map
		// isinya accesstoken dan id token, dan akan di unmarshal kedalam variable googleOauthToken
		err = json.Unmarshal([]byte(googleToken), &googleOauthToken)
		if err != nil {
			log.Err(err).Msg("failed unmarshal google token")
			return nil, _error.ErrStringDefault(http.StatusForbidden)
		}
	} else {
		// kita decrypt request body dari user web make aes cbc
		googleCode, err := encryption.DecryptStringCBC(req.Token, config.AesCBC, config.AesCBCIV)
		if err != nil {
			return nil, _error.ErrStringDefault(http.StatusForbidden)
		}

		// ini isinya string code token dari fe web
		// yang nnti nya akan dipakai untuk generate acces token dan id token untuk memasukann.
		// kedalam variable googleOauthToken
		googleOauthToken, err = o.oauth2Repo.GetGoogleOauthToken(googleCode)
		if err != nil {
			return nil, _error.ErrStringDefault(http.StatusForbidden)
		}
	}

	// GetGoogleOauthUser get google user menggunakan access token dan id token
	googleUser, err := o.oauth2Repo.GetGoogleOauthUser(googleOauthToken)
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusForbidden)
	}

	// OpenConn open connection database dari userrepo
	// defer untuk melakukan close connection
	err = o.userRepo.OpenConn(ctx)
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer o.userRepo.CloseConn()

	// kita check apakah email sudah terdaftar apa belum
	exist, err := o.userRepo.CheckUserByEmail(ctx, googleUser.Email)
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	// kita convert dari hasil GetGoogleOauthUser ke dalam response
	// serta memasukan status apakah user sudah terdaftar atau belum
	googleUserResp := converter.GoogleClaimModelToResp(googleUser, exist)
	return googleUserResp, nil
}
