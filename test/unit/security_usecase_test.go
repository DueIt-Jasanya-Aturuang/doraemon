package unit

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/mocks"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/pkg/_usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/pkg/helper"
)

func TestSecurityUsecaseJwtValidateAT(t *testing.T) {
	infra.EnvInit()
	userRepo := &mocks.FakeUserSqlRepo{}
	securityRepo := &mocks.FakeSecuritySqlRepo{}
	securityUsecase := _usecase.NewSecurityUsecaseImpl(userRepo, securityRepo)

	userID := "userID_1"
	var jwtModel *model.Jwt
	jwtModel = jwtModel.AccessTokenDefault(userID)
	accessToken, err := helper.GenerateJwtHS256(jwtModel)
	assert.NoError(t, err)

	jwtModel = jwtModel.RefreshTokenDefault(userID, true)
	refreshToken, err := helper.GenerateJwtHS256(jwtModel)
	assert.NoError(t, err)

	req := &dto.JwtTokenReq{
		AppId:         "appID_1",
		Authorization: accessToken,
		UserId:        userID,
	}

	tokenModel := &model.Token{
		UserID:       userID,
		AppID:        "appID_1",
		RefreshToken: refreshToken,
		AcceesToken:  accessToken,
		RememberMe:   false,
	}

	t.Run("SUCCESS", func(t *testing.T) {
		_ = securityRepo.OpenConn(context.TODO())
		securityRepo.OpenConnReturns(nil)
		defer securityRepo.CloseConn()

		_, _ = securityRepo.GetTokenByAT(context.TODO(), accessToken)
		securityRepo.GetTokenByATReturns(tokenModel, nil)

		userRepo.CheckActivasiUserByID(context.TODO(), "id")
		userRepo.CheckActivasiUserByIDReturns(true, nil)
		exp, err := securityUsecase.JwtValidateAT(context.TODO(), req, "not activasi")
		assert.NoError(t, err)
		assert.Equal(t, false, exp)
	})

	t.Run("SUCCESS_with-check-activasi", func(t *testing.T) {
		_ = securityRepo.OpenConn(context.TODO())
		securityRepo.OpenConnReturns(nil)
		defer securityRepo.CloseConn()

		_, _ = securityRepo.GetTokenByAT(context.TODO(), accessToken)
		securityRepo.GetTokenByATReturns(tokenModel, nil)

		_, _ = userRepo.CheckActivasiUserByID(context.TODO(), userID)
		userRepo.CheckActivasiUserByIDReturns(true, nil)

		exp, err := securityUsecase.JwtValidateAT(context.TODO(), req, "/activasi-account")
		assert.NoError(t, err)
		assert.Equal(t, false, exp)
	})

	t.Run("ERROR_token-nil", func(t *testing.T) {
		_ = securityRepo.OpenConn(context.TODO())
		securityRepo.OpenConnReturns(nil)
		defer securityRepo.CloseConn()

		_, _ = securityRepo.GetTokenByAT(context.TODO(), accessToken)
		securityRepo.GetTokenByATReturns(nil, sql.ErrNoRows)

		_ = securityRepo.DeleteAllTokenByUserID(context.TODO(), req.UserId)
		securityRepo.DeleteAllTokenByUserIDReturns(nil)

		exp, err := securityUsecase.JwtValidateAT(context.TODO(), req, "/activasi-account")
		assert.NotNil(t, err)
		assert.Equal(t, false, exp)
	})

	t.Run("ERROR_exp-token", func(t *testing.T) {
		req.Authorization = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTM5OTMyOTYsInN1YiI6IjEyMyJ9.d4eTOPwf54uVoMKghgpuC9BcJKhUgt-KI66Ncfcssew"
		_ = securityRepo.OpenConn(context.TODO())
		securityRepo.OpenConnReturns(nil)
		defer securityRepo.CloseConn()

		_, _ = securityRepo.GetTokenByAT(context.TODO(), accessToken)
		securityRepo.GetTokenByATReturns(tokenModel, nil)

		_, _ = userRepo.CheckActivasiUserByID(context.TODO(), userID)
		userRepo.CheckActivasiUserByIDReturns(false, nil)

		exp, err := securityUsecase.JwtValidateAT(context.TODO(), req, "/activasi-account")
		assert.Nil(t, err)
		assert.Equal(t, true, exp)
	})

	t.Run("ERROR_invalid-token", func(t *testing.T) {
		_ = securityRepo.OpenConn(context.TODO())
		securityRepo.OpenConnReturns(nil)
		defer securityRepo.CloseConn()

		req.Authorization = "invalid"
		exp, err := securityUsecase.JwtValidateAT(context.TODO(), req, "/activasi-account")
		assert.Error(t, err)
		assert.Equal(t, false, exp)
		var errHTTP *model.ErrResponseHTTP
		assert.Equal(t, true, errors.As(err, &errHTTP))
		assert.Equal(t, 401, errHTTP.Code)
	})

	t.Run("ERROR_with-check-activasi", func(t *testing.T) {
		req.Authorization = accessToken
		_ = securityRepo.OpenConn(context.TODO())
		securityRepo.OpenConnReturns(nil)
		defer securityRepo.CloseConn()

		_, _ = securityRepo.GetTokenByAT(context.TODO(), accessToken)
		securityRepo.GetTokenByATReturns(tokenModel, nil)

		_, _ = userRepo.CheckActivasiUserByID(context.TODO(), userID)
		userRepo.CheckActivasiUserByIDReturns(false, nil)

		exp, err := securityUsecase.JwtValidateAT(context.TODO(), req, "/activasi-accunt")
		assert.Error(t, err)
		assert.Equal(t, false, exp)
		var errHTTP *model.ErrResponseHTTP
		assert.Equal(t, true, errors.As(err, &errHTTP))
		assert.Equal(t, 403, errHTTP.Code)
	})
}

func TestSecurityUsecaseJwtGenerateRTAT(t *testing.T) {
	infra.EnvInit()
	userRepo := &mocks.FakeUserSqlRepo{}
	securityRepo := &mocks.FakeSecuritySqlRepo{}
	securityUsecase := _usecase.NewSecurityUsecaseImpl(userRepo, securityRepo)

	userID := "userID_1"
	var jwtModel *model.Jwt
	jwtModel = jwtModel.AccessTokenDefault(userID)
	accessToken, err := helper.GenerateJwtHS256(jwtModel)
	assert.NoError(t, err)

	jwtModel = jwtModel.RefreshTokenDefault(userID, true)
	refreshToken, err := helper.GenerateJwtHS256(jwtModel)
	assert.NoError(t, err)

	req := &dto.JwtTokenReq{
		AppId:         "appID_1",
		Authorization: accessToken,
		UserId:        userID,
	}

	tokenModel := &model.Token{
		UserID:       userID,
		AppID:        "appID_1",
		RefreshToken: refreshToken,
		AcceesToken:  accessToken,
		RememberMe:   false,
	}

	t.Run("SUCCESS", func(t *testing.T) {
		_ = securityRepo.OpenConn(context.TODO())
		securityRepo.OpenConnReturns(nil)
		defer securityRepo.CloseConn()

		_ = securityRepo.StartTx(context.TODO(), &sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
			ReadOnly:  false,
		})
		securityRepo.StartTxReturns(nil)

		_, _ = securityRepo.GetTokenByAT(context.TODO(), accessToken)
		securityRepo.GetTokenByATReturns(tokenModel, nil)

		_ = securityRepo.DeleteToken(context.TODO(), 1, req.UserId)
		securityRepo.DeleteTokenReturns(nil)

		var jwtModel *model.Jwt
		jwtModel = jwtModel.AccessTokenDefault(userID)
		accessToken, err := helper.GenerateJwtHS256(jwtModel)
		assert.NoError(t, err)

		jwtModel = jwtModel.RefreshTokenDefault(userID, true)
		refreshToken, err := helper.GenerateJwtHS256(jwtModel)
		assert.NoError(t, err)

		_ = securityRepo.UpdateToken(context.TODO(), 1, refreshToken, accessToken)
		securityRepo.UpdateTokenReturns(nil)

		tokenResp, err := securityUsecase.JwtGenerateRTAT(context.TODO(), req)
		_ = securityRepo.EndTx(err)
		securityRepo.EndTxReturns(err)

		assert.NoError(t, err)
		assert.NotNil(t, tokenResp)
	})

	t.Run("ERROR_token-nil", func(t *testing.T) {
		_ = securityRepo.OpenConn(context.TODO())
		securityRepo.OpenConnReturns(nil)
		defer securityRepo.CloseConn()

		_ = securityRepo.StartTx(context.TODO(), &sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
			ReadOnly:  false,
		})
		securityRepo.StartTxReturns(nil)

		_, _ = securityRepo.GetTokenByAT(context.TODO(), accessToken)
		securityRepo.GetTokenByATReturns(nil, sql.ErrNoRows)

		_ = securityRepo.DeleteAllTokenByUserID(context.TODO(), req.UserId)
		securityRepo.DeleteAllTokenByUserIDReturns(nil)

		tokenResp, err := securityUsecase.JwtGenerateRTAT(context.TODO(), req)
		_ = securityRepo.EndTx(err)
		securityRepo.EndTxReturns(err)

		assert.Error(t, err)
		assert.Nil(t, tokenResp)
		var errHTTP *model.ErrResponseHTTP
		assert.Equal(t, true, errors.As(err, &errHTTP))
		assert.Equal(t, 401, errHTTP.Code)
	})

	t.Run("ERROR_token-exp", func(t *testing.T) {
		// nunggu 24 jam sampe refresh token expired
		tokenModel.RefreshToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTQwMDEwNDEsInN1YiI6IjEyMyJ9.JvK6TeNzVBQQEoVkoLKGyzTwM-9R5bHiyX2_P-b6Ti0"
		_ = securityRepo.OpenConn(context.TODO())
		securityRepo.OpenConnReturns(nil)
		defer securityRepo.CloseConn()

		_ = securityRepo.StartTx(context.TODO(), &sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
			ReadOnly:  false,
		})
		securityRepo.StartTxReturns(nil)

		_, _ = securityRepo.GetTokenByAT(context.TODO(), accessToken)
		securityRepo.GetTokenByATReturns(tokenModel, nil)

		_ = securityRepo.DeleteToken(context.TODO(), 1, req.UserId)
		securityRepo.DeleteTokenReturns(nil)

		tokenResp, err := securityUsecase.JwtGenerateRTAT(context.TODO(), req)
		_ = securityRepo.EndTx(err)
		securityRepo.EndTxReturns(err)

		assert.Error(t, err)
		assert.Nil(t, tokenResp)
		var errHTTP *model.ErrResponseHTTP
		assert.Equal(t, true, errors.As(err, &errHTTP))
		assert.Equal(t, 401, errHTTP.Code)
	})
}

func TestSecurityUsecaseJwtRegistredRTAT(t *testing.T) {
	infra.EnvInit()
	userRepo := &mocks.FakeUserSqlRepo{}
	securityRepo := &mocks.FakeSecuritySqlRepo{}
	securityUsecase := _usecase.NewSecurityUsecaseImpl(userRepo, securityRepo)

	userID := "userID_1"
	var jwtModel *model.Jwt
	jwtModel = jwtModel.AccessTokenDefault(userID)
	accessToken, err := helper.GenerateJwtHS256(jwtModel)
	assert.NoError(t, err)

	jwtModel = jwtModel.RefreshTokenDefault(userID, true)
	refreshToken, err := helper.GenerateJwtHS256(jwtModel)
	assert.NoError(t, err)

	req := &dto.JwtRegisteredTokenReq{
		AppId:      "appID_1",
		UserId:     userID,
		RememberMe: false,
	}

	tokenModel := &model.Token{
		UserID:       userID,
		AppID:        "appID_1",
		RefreshToken: refreshToken,
		AcceesToken:  accessToken,
		RememberMe:   false,
	}

	t.Run("SUCCESS", func(t *testing.T) {
		_ = securityRepo.OpenConn(context.TODO())
		securityRepo.OpenConnReturns(nil)
		defer securityRepo.CloseConn()

		_ = securityRepo.StartTx(context.TODO(), &sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
			ReadOnly:  false,
		})
		securityRepo.StartTxReturns(nil)

		_ = securityRepo.CreateToken(context.TODO(), tokenModel)
		securityRepo.CreateTokenReturns(nil)

		tokenResp, err := securityUsecase.JwtRegistredRTAT(context.TODO(), req)
		assert.NoError(t, err)
		assert.NotNil(t, tokenResp)
	})

}

func TestSecurityUsecaseLogout(t *testing.T) {
	infra.EnvInit()
	userRepo := &mocks.FakeUserSqlRepo{}
	securityRepo := &mocks.FakeSecuritySqlRepo{}
	securityUsecase := _usecase.NewSecurityUsecaseImpl(userRepo, securityRepo)

	userID := "userID_1"

	var jwtModel *model.Jwt
	jwtModel = jwtModel.AccessTokenDefault(userID)
	accessToken, err := helper.GenerateJwtHS256(jwtModel)
	assert.NoError(t, err)

	jwtModel = jwtModel.RefreshTokenDefault(userID, true)
	refreshToken, err := helper.GenerateJwtHS256(jwtModel)
	assert.NoError(t, err)

	tokenModel := &model.Token{
		UserID:       userID,
		AppID:        "appID_1",
		RefreshToken: refreshToken,
		AcceesToken:  accessToken,
		RememberMe:   false,
	}

	t.Run("SUCCESS", func(t *testing.T) {
		err = securityRepo.OpenConn(context.TODO())
		securityRepo.OpenConnReturns(nil)
		assert.NoError(t, err)
		defer securityRepo.CloseConn()

		err = securityRepo.StartTx(context.TODO(), &sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
			ReadOnly:  false,
		})
		assert.NoError(t, err)
		securityRepo.StartTxReturns(nil)
		defer func() {
			errEndTx := securityRepo.EndTx(err)
			assert.NoError(t, errEndTx)
		}()

		_, _ = securityRepo.GetTokenByAT(context.TODO(), accessToken)
		securityRepo.GetTokenByATReturns(tokenModel, nil)

		err = securityUsecase.Logout(context.TODO(), &dto.LogoutReq{
			Token:  accessToken,
			UserID: userID,
		})
		assert.NoError(t, err)
	})

	t.Run("ERROR_token-nil", func(t *testing.T) {
		err = securityRepo.OpenConn(context.TODO())
		securityRepo.OpenConnReturns(nil)
		assert.NoError(t, err)
		defer securityRepo.CloseConn()

		err = securityRepo.StartTx(context.TODO(), &sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
			ReadOnly:  false,
		})
		assert.NoError(t, err)
		securityRepo.StartTxReturns(nil)
		defer func() {
			errEndTx := securityRepo.EndTx(err)
			assert.NoError(t, errEndTx)
		}()

		_, _ = securityRepo.GetTokenByAT(context.TODO(), accessToken)
		securityRepo.GetTokenByATReturns(nil, sql.ErrNoRows)

		err = securityUsecase.Logout(context.TODO(), &dto.LogoutReq{
			Token:  accessToken,
			UserID: userID,
		})
		assert.NoError(t, err)
	})
}
