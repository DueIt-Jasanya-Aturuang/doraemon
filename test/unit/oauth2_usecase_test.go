package unit

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/mocks"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/pkg/_usecase"
)

func TestOauth2UsecaseGoogleClaimUser(t *testing.T) {
	infra.EnvInit()
	userRepo := &mocks.FakeUserSqlRepo{}
	oauth2Repo := &mocks.FakeOauth2ProviderRepo{}

	oauth2Usecase := _usecase.NewOauth2UsecaseImpl(userRepo, oauth2Repo)

	googleTokenModel := &model.GoogleOauth2Token{
		AccessToken: "access_token",
		IDToken:     "id_token",
	}

	googleUserModel := &model.GoogleOauth2User{
		ID:            "googleID_1",
		Email:         "ibanrama29@gmail.com",
		VerifiedEmail: true,
		Name:          "rama",
		GivenName:     "ibanrama",
		FamilyName:    "rama",
		Image:         "https://googleimage",
		Locale:        "indonesia",
	}

	t.Run("SUCCESS_web", func(t *testing.T) {
		req := &dto.LoginGoogleReq{
			Token:  "GYaFpgGgx8HvLd+elVnfFA==",
			Device: "web",
		}

		_, _ = oauth2Repo.GetGoogleOauthToken(req.Token)
		oauth2Repo.GetGoogleOauthTokenReturns(googleTokenModel, nil)

		_, _ = oauth2Repo.GetGoogleOauthUser(googleTokenModel)
		oauth2Repo.GetGoogleOauthUserReturns(googleUserModel, nil)

		_ = userRepo.OpenConn(context.TODO())
		userRepo.OpenConnReturns(nil)
		defer userRepo.CloseConn()

		_, _ = userRepo.CheckUserByEmail(context.TODO(), googleUserModel.Email)
		userRepo.CheckUserByEmailReturns(false, nil)

		googleResp, err := oauth2Usecase.GoogleClaimUser(context.TODO(), req)
		assert.NoError(t, err)
		assert.NotNil(t, googleResp)
	})

	t.Run("SUCCESS_mobile", func(t *testing.T) {
		req := &dto.LoginGoogleReq{
			Token:  "62idN0WO0Ym9cqcKnCPus8HGj1RE8SZ54IgGFLyRgw6BPjtfsV47NnBfKV3jSZEi8f8mhk8Olu+OslX/2a2Oi1S+Jlwh6zRezSSMPqwjRug=",
			Device: "mobile",
		}

		_, _ = oauth2Repo.GetGoogleOauthUser(googleTokenModel)
		oauth2Repo.GetGoogleOauthUserReturns(googleUserModel, nil)

		_ = userRepo.OpenConn(context.TODO())
		userRepo.OpenConnReturns(nil)
		defer userRepo.CloseConn()

		_, _ = userRepo.CheckUserByEmail(context.TODO(), googleUserModel.Email)
		userRepo.CheckUserByEmailReturns(false, nil)

		googleResp, err := oauth2Usecase.GoogleClaimUser(context.TODO(), req)
		assert.NoError(t, err)
		assert.NotNil(t, googleResp)
	})

	t.Run("ERROR_invalid-token-req-mobile", func(t *testing.T) {
		req := &dto.LoginGoogleReq{
			Token:  "ccLHBNyHs5gA+K3hBAHusg==",
			Device: "mobile",
		}

		_, _ = oauth2Repo.GetGoogleOauthUser(googleTokenModel)
		oauth2Repo.GetGoogleOauthUserReturns(googleUserModel, nil)

		_ = userRepo.OpenConn(context.TODO())
		userRepo.OpenConnReturns(nil)
		defer userRepo.CloseConn()

		_, _ = userRepo.CheckUserByEmail(context.TODO(), googleUserModel.Email)
		userRepo.CheckUserByEmailReturns(false, nil)

		googleResp, err := oauth2Usecase.GoogleClaimUser(context.TODO(), req)
		assert.Error(t, err)
		assert.Nil(t, googleResp)
		var errHTTP *model.ErrResponseHTTP
		assert.Equal(t, true, errors.As(err, &errHTTP))
	})

	t.Run("ERROR_invalid-token-req-web", func(t *testing.T) {
		req := &dto.LoginGoogleReq{
			Token:  "ccLHBNyHs5gA+K3hBAHusg==",
			Device: "web",
		}

		_, _ = oauth2Repo.GetGoogleOauthToken(req.Token)
		oauth2Repo.GetGoogleOauthTokenReturns(googleTokenModel, nil)

		_, _ = oauth2Repo.GetGoogleOauthUser(googleTokenModel)
		oauth2Repo.GetGoogleOauthUserReturns(googleUserModel, nil)

		_ = userRepo.OpenConn(context.TODO())
		userRepo.OpenConnReturns(nil)
		defer userRepo.CloseConn()

		_, _ = userRepo.CheckUserByEmail(context.TODO(), googleUserModel.Email)
		userRepo.CheckUserByEmailReturns(false, nil)

		googleResp, err := oauth2Usecase.GoogleClaimUser(context.TODO(), req)
		assert.Error(t, err)
		assert.Nil(t, googleResp)
		var errHTTP *model.ErrResponseHTTP
		assert.Equal(t, true, errors.As(err, &errHTTP))
	})

	t.Run("ERROR_error-get-user-google-mobile", func(t *testing.T) {
		req := &dto.LoginGoogleReq{
			Token:  "62idN0WO0Ym9cqcKnCPus8HGj1RE8SZ54IgGFLyRgw6BPjtfsV47NnBfKV3jSZEi8f8mhk8Olu+OslX/2a2Oi1S+Jlwh6zRezSSMPqwjRug=",
			Device: "mobile",
		}

		_, _ = oauth2Repo.GetGoogleOauthUser(googleTokenModel)
		oauth2Repo.GetGoogleOauthUserReturns(nil, errors.New("forbidden"))

		_ = userRepo.OpenConn(context.TODO())
		userRepo.OpenConnReturns(nil)
		defer userRepo.CloseConn()

		_, _ = userRepo.CheckUserByEmail(context.TODO(), googleUserModel.Email)
		userRepo.CheckUserByEmailReturns(false, nil)

		googleResp, err := oauth2Usecase.GoogleClaimUser(context.TODO(), req)
		assert.Error(t, err)
		assert.Nil(t, googleResp)
		var errHTTP *model.ErrResponseHTTP
		assert.Equal(t, true, errors.As(err, &errHTTP))
	})

	t.Run("ERROR_error-get-token-google-web", func(t *testing.T) {
		req := &dto.LoginGoogleReq{
			Token:  "GYaFpgGgx8HvLd+elVnfFA==",
			Device: "web",
		}

		_, _ = oauth2Repo.GetGoogleOauthToken(req.Token)
		oauth2Repo.GetGoogleOauthTokenReturns(nil, errors.New("forbidden"))

		_, _ = oauth2Repo.GetGoogleOauthUser(googleTokenModel)
		oauth2Repo.GetGoogleOauthUserReturns(nil, nil)

		_ = userRepo.OpenConn(context.TODO())
		userRepo.OpenConnReturns(nil)
		defer userRepo.CloseConn()

		_, _ = userRepo.CheckUserByEmail(context.TODO(), googleUserModel.Email)
		userRepo.CheckUserByEmailReturns(false, nil)

		googleResp, err := oauth2Usecase.GoogleClaimUser(context.TODO(), req)
		assert.Error(t, err)
		assert.Nil(t, googleResp)
		var errHTTP *model.ErrResponseHTTP
		assert.Equal(t, true, errors.As(err, &errHTTP))
	})

	t.Run("ERROR_error-get-user-google-web", func(t *testing.T) {
		req := &dto.LoginGoogleReq{
			Token:  "GYaFpgGgx8HvLd+elVnfFA==",
			Device: "web",
		}

		_, _ = oauth2Repo.GetGoogleOauthToken(req.Token)
		oauth2Repo.GetGoogleOauthTokenReturns(googleTokenModel, nil)

		_, _ = oauth2Repo.GetGoogleOauthUser(googleTokenModel)
		oauth2Repo.GetGoogleOauthUserReturns(nil, errors.New("forbidden"))

		_ = userRepo.OpenConn(context.TODO())
		userRepo.OpenConnReturns(nil)
		defer userRepo.CloseConn()

		_, _ = userRepo.CheckUserByEmail(context.TODO(), googleUserModel.Email)
		userRepo.CheckUserByEmailReturns(false, nil)

		googleResp, err := oauth2Usecase.GoogleClaimUser(context.TODO(), req)
		assert.Error(t, err)
		assert.Nil(t, googleResp)
		var errHTTP *model.ErrResponseHTTP
		assert.Equal(t, true, errors.As(err, &errHTTP))
	})
}
