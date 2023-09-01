package unit

import (
	"errors"
	"testing"
	"time"
	
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infrastructures/config"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/helper"
)

func TestJwtHelper(t *testing.T) {
	config.EnvInit()
	var accessTokenTrue string
	var refreshTokenTrue string
	var forgotPasswordToken string
	var accessTokenFalse string
	var refreshTokenFalse string
	oldToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTM1OTE5ODIsInN1YiI6ImVhYzM5YjBmLTRjNzctNDZiZS1iMDEyLWU1ODYyMzhkOGMyMzpmYWxzZTpmb3Jnb3RfcGFzc3dvcmQifQ.02K96TTbMnDEsJbqinUEs-xpLIUytZryDwDJJA974h4"
	userID := "123"
	t.Run("SUCCESS_false", func(t *testing.T) {
		tUUID := uuid.NewV4().String()

		var jwtModel *model.Jwt

		jwtModel = jwtModel.AccessTokenDefault(tUUID, userID, false)
		accessToken, err := helper.GenerateJwtHS256(jwtModel)
		assert.NoError(t, err)
		t.Log(accessToken)
		accessTokenFalse = accessToken

		jwtModel = jwtModel.RefreshTokenDefault(tUUID, userID, false)
		refreshTokenRes, err := helper.GenerateJwtHS256(jwtModel)
		assert.NoError(t, err)
		refreshTokenFalse = refreshTokenRes
		t.Log(refreshTokenRes)

		jwtModel = jwtModel.ForgotPasswordTokenDefault(tUUID, userID)
		forgotPasswordTokenRes, err := helper.GenerateJwtHS256(jwtModel)
		assert.NoError(t, err)
		t.Log(forgotPasswordToken)
		forgotPasswordToken = forgotPasswordTokenRes
	})

	t.Run("SUCCESS_true", func(t *testing.T) {
		tUUID := uuid.NewV4().String()

		var jwtModel *model.Jwt

		jwtModel = jwtModel.AccessTokenDefault(tUUID, userID, true)
		accessToken, err := helper.GenerateJwtHS256(jwtModel)
		assert.NoError(t, err)
		t.Log(accessToken)
		accessTokenTrue = accessToken

		jwtModel = jwtModel.RefreshTokenDefault(tUUID, userID, true)
		refreshToken, err := helper.GenerateJwtHS256(jwtModel)
		assert.NoError(t, err)
		t.Log(refreshToken)
		refreshTokenTrue = refreshToken

	})

	t.Run("SUCCESS_CLAIM", func(t *testing.T) {
		atClaims, err := helper.ClaimsJwtHS256(accessTokenTrue, config.AccessTokenKeyHS)
		assert.NoError(t, err)
		t.Log(atClaims)

		rtClaims, err := helper.ClaimsJwtHS256(refreshTokenTrue, config.RefreshTokenKeyHS)
		assert.NoError(t, err)
		t.Log(rtClaims)

		atClaimsFalse, err := helper.ClaimsJwtHS256(accessTokenFalse, config.AccessTokenKeyHS)
		assert.NoError(t, err)
		t.Log(atClaimsFalse)

		rtClaimFalse, err := helper.ClaimsJwtHS256(refreshTokenFalse, config.RefreshTokenKeyHS)
		assert.NoError(t, err)
		t.Log(rtClaimFalse)

		fpClaim, err := helper.ClaimsJwtHS256(forgotPasswordToken, config.DefaultKey)
		assert.NoError(t, err)
		t.Log(fpClaim)

		rtTrueClaim := rtClaims["exp"].(float64)
		unix := time.Unix(int64(rtTrueClaim), 0)
		t.Log(unix)
	})

	t.Run("ERROR_CLAIM", func(t *testing.T) {
		_, err := helper.ClaimsJwtHS256(accessTokenTrue, config.RefreshTokenKeyHS)
		assert.Error(t, err)
		t.Log(err)

		_, err = helper.ClaimsJwtHS256(oldToken, config.DefaultKey)
		assert.Error(t, err)
		con := errors.Is(err, jwt.ErrTokenExpired)
		assert.Equal(t, true, con)
		t.Log(err)
	})
}
