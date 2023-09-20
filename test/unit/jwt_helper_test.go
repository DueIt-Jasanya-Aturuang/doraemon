package unit

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/helper"
)

func TestJwtHelper(t *testing.T) {
	infra.EnvInit()
	var accessTokenTrue string
	var refreshTokenTrue string
	var forgotPasswordToken string
	var accessTokenFalse string
	var refreshTokenFalse string
	oldToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTM1OTE5ODIsInN1YiI6ImVhYzM5YjBmLTRjNzctNDZiZS1iMDEyLWU1ODYyMzhkOGMyMzpmYWxzZTpmb3Jnb3RfcGFzc3dvcmQifQ.02K96TTbMnDEsJbqinUEs-xpLIUytZryDwDJJA974h4"
	userID := "123"
	t.Run("SUCCESS_false", func(t *testing.T) {
		var jwtModel *model.Jwt

		jwtModel = jwtModel.AccessTokenDefault(userID)
		accessToken, err := helper.GenerateJwtHS256(jwtModel)
		assert.NoError(t, err)
		t.Log(accessToken)
		accessTokenFalse = accessToken

		jwtModel = jwtModel.RefreshTokenDefault(userID, false)
		refreshTokenRes, err := helper.GenerateJwtHS256(jwtModel)
		assert.NoError(t, err)
		refreshTokenFalse = refreshTokenRes
		t.Log(refreshTokenRes)

		jwtModel = jwtModel.ForgotPasswordTokenDefault(userID)
		forgotPasswordTokenRes, err := helper.GenerateJwtHS256(jwtModel)
		assert.NoError(t, err)
		t.Log(forgotPasswordToken)
		forgotPasswordToken = forgotPasswordTokenRes
	})

	t.Run("SUCCESS_true", func(t *testing.T) {
		var jwtModel *model.Jwt

		jwtModel = jwtModel.AccessTokenDefault(userID)
		accessToken, err := helper.GenerateJwtHS256(jwtModel)
		assert.NoError(t, err)
		t.Log(accessToken)
		accessTokenTrue = accessToken

		jwtModel = jwtModel.RefreshTokenDefault(userID, true)
		refreshToken, err := helper.GenerateJwtHS256(jwtModel)
		assert.NoError(t, err)
		t.Log(refreshToken)
		refreshTokenTrue = refreshToken

	})

	t.Run("SUCCESS_CLAIM", func(t *testing.T) {
		atClaims, err := helper.ClaimsJwtHS256(accessTokenTrue, infra.AccessTokenKeyHS)
		assert.NoError(t, err)
		t.Log(atClaims)

		rtClaims, err := helper.ClaimsJwtHS256(refreshTokenTrue, infra.RefreshTokenKeyHS)
		assert.NoError(t, err)
		t.Log(rtClaims)

		atClaimsFalse, err := helper.ClaimsJwtHS256(accessTokenFalse, infra.AccessTokenKeyHS)
		assert.NoError(t, err)
		t.Log(atClaimsFalse)

		rtClaimFalse, err := helper.ClaimsJwtHS256(refreshTokenFalse, infra.RefreshTokenKeyHS)
		assert.NoError(t, err)
		t.Log(rtClaimFalse)

		fpClaim, err := helper.ClaimsJwtHS256(forgotPasswordToken, infra.DefaultKey)
		assert.NoError(t, err)
		t.Log(fpClaim)

		rtTrueClaim := rtClaims["exp"].(float64)
		unix := time.Unix(int64(rtTrueClaim), 0)
		t.Log(unix)
	})

	t.Run("ERROR_CLAIM", func(t *testing.T) {
		_, err := helper.ClaimsJwtHS256(accessTokenTrue, infra.RefreshTokenKeyHS)
		assert.Error(t, err)
		t.Log(err)

		claim, err := helper.ClaimsJwtHS256(oldToken, infra.DefaultKey)
		assert.Error(t, err)
		t.Log(err)
		t.Log(claim)
		con := errors.Is(err, jwt.ErrTokenExpired)
		assert.Equal(t, true, con)
	})
}
