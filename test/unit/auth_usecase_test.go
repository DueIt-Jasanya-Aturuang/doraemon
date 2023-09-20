package unit

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/mocks"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/_usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/converter"
)

func TestAuthUsecaseLogin(t *testing.T) {
	userRepo := &mocks.FakeUserSqlRepo{}
	accessRepo := &mocks.FakeAccessSqlRepo{}

	accountApi := &mocks.FakeAccountApiRepo{}

	authUsecase := _usecase.NewAuthUsecaseImpl(userRepo, accessRepo, accountApi)

	user := &model.User{
		ID:              "userID_1",
		FullName:        "rama",
		Gender:          "undefined",
		Image:           "/files/user-image/public/default-male.png",
		Username:        "iban.rama",
		Email:           "ibanrama29@gmail.com",
		Password:        "$2a$14$XmVO9LPkf1aDBUF8M/3jqOt8TPizQC1j5p1Hf.VhHHmA0kviudEm6",
		PhoneNumber:     sql.NullString{},
		EmailVerifiedAt: false,
		CreatedAt:       0,
		CreatedBy:       "userID_1",
		UpdatedAt:       0,
		UpdatedBy:       sql.NullString{},
		DeletedAt:       sql.NullInt64{},
		DeletedBy:       sql.NullString{},
	}

	profile := &model.Profile{
		Code:      200,
		ProfileID: "profileID_1",
		Quote:     "null",
		Profesi:   "null",
	}

	t.Run("SUCCESS_email", func(t *testing.T) {
		req := &dto.LoginReq{
			EmailOrUsername: "ibanrama29@gmail.com",
			Password:        "rama123",
			RememberMe:      true,
			Oauth2:          false,
		}
		err := userRepo.OpenConn(context.TODO())
		userRepo.OpenConnReturns(nil)
		assert.NoError(t, err)
		defer userRepo.CloseConn()

		_, err = userRepo.GetUserByEmailOrUsername(context.TODO(), req.EmailOrUsername)
		assert.NoError(t, err)
		userRepo.GetUserByEmailOrUsernameReturns(user, nil)

		_, err = accountApi.GetProfileByUserID(user.ID)
		assert.NoError(t, err)
		accountApi.GetProfileByUserIDReturns(profile, nil)

		userResp, profileResp, err := authUsecase.Login(context.TODO(), req)
		assert.NoError(t, err)
		assert.NotNil(t, userResp)
		assert.NotNil(t, profileResp)
		assert.Equal(t, profile.ProfileID, profileResp.ProfileID)
		assert.Equal(t, user.ID, user.ID)
	})

	t.Run("SUCCESS_username", func(t *testing.T) {
		req := &dto.LoginReq{
			EmailOrUsername: "iban.rama",
			Password:        "rama123",
			RememberMe:      true,
			Oauth2:          false,
		}
		err := userRepo.OpenConn(context.TODO())
		userRepo.OpenConnReturns(nil)
		assert.NoError(t, err)
		defer userRepo.CloseConn()

		_, err = userRepo.GetUserByEmailOrUsername(context.TODO(), req.EmailOrUsername)
		assert.NoError(t, err)
		userRepo.GetUserByEmailOrUsernameReturns(user, nil)

		_, err = accountApi.GetProfileByUserID(user.ID)
		assert.NoError(t, err)
		accountApi.GetProfileByUserIDReturns(profile, nil)

		userResp, profileResp, err := authUsecase.Login(context.TODO(), req)
		assert.NoError(t, err)
		assert.NotNil(t, userResp)
		assert.NotNil(t, profileResp)
		assert.Equal(t, profile.ProfileID, profileResp.ProfileID)
		assert.Equal(t, user.ID, user.ID)
	})

	t.Run("SUCCESS_oauth2", func(t *testing.T) {
		req := &dto.LoginReq{
			EmailOrUsername: "iban.rama",
			Password:        "invalid password",
			RememberMe:      true,
			Oauth2:          true,
		}
		err := userRepo.OpenConn(context.TODO())
		userRepo.OpenConnReturns(nil)
		assert.NoError(t, err)
		defer userRepo.CloseConn()

		_, err = userRepo.GetUserByEmailOrUsername(context.TODO(), req.EmailOrUsername)
		assert.NoError(t, err)
		userRepo.GetUserByEmailOrUsernameReturns(user, nil)

		_, err = accountApi.GetProfileByUserID(user.ID)
		assert.NoError(t, err)
		accountApi.GetProfileByUserIDReturns(profile, nil)

		userResp, profileResp, err := authUsecase.Login(context.TODO(), req)
		assert.NoError(t, err)
		assert.NotNil(t, userResp)
		assert.NotNil(t, profileResp)
		assert.Equal(t, profile.ProfileID, profileResp.ProfileID)
		assert.Equal(t, user.ID, user.ID)
	})

	t.Run("ERROR_invalid-password", func(t *testing.T) {
		req := &dto.LoginReq{
			EmailOrUsername: "iban.rama",
			Password:        "rama1234",
			RememberMe:      true,
			Oauth2:          false,
		}
		err := userRepo.OpenConn(context.TODO())
		userRepo.OpenConnReturns(nil)
		assert.NoError(t, err)
		defer userRepo.CloseConn()

		_, err = userRepo.GetUserByEmailOrUsername(context.TODO(), req.EmailOrUsername)
		assert.NoError(t, err)
		userRepo.GetUserByEmailOrUsernameReturns(user, nil)

		_, err = accountApi.GetProfileByUserID(user.ID)
		assert.NoError(t, err)
		accountApi.GetProfileByUserIDReturns(profile, nil)

		userResp, profileResp, err := authUsecase.Login(context.TODO(), req)
		assert.Error(t, err)
		var errHTTP *model.ErrResponseHTTP
		assert.Equal(t, true, errors.As(err, &errHTTP))
		assert.Equal(t, _error.BadLogin(), errHTTP)
		assert.Nil(t, userResp)
		assert.Nil(t, profileResp)
	})

	t.Run("ERROR_empty-user", func(t *testing.T) {
		req := &dto.LoginReq{
			EmailOrUsername: "iban.rama",
			Password:        "rama1234",
			RememberMe:      true,
			Oauth2:          false,
		}
		err := userRepo.OpenConn(context.TODO())
		userRepo.OpenConnReturns(nil)
		assert.NoError(t, err)
		defer userRepo.CloseConn()

		_, err = userRepo.GetUserByEmailOrUsername(context.TODO(), req.EmailOrUsername)
		assert.NoError(t, err)
		userRepo.GetUserByEmailOrUsernameReturns(nil, sql.ErrNoRows)

		_, err = accountApi.GetProfileByUserID(user.ID)
		assert.NoError(t, err)
		accountApi.GetProfileByUserIDReturns(profile, nil)

		userResp, profileResp, err := authUsecase.Login(context.TODO(), req)
		assert.Error(t, err)
		var errHTTP *model.ErrResponseHTTP
		assert.Equal(t, true, errors.As(err, &errHTTP))
		assert.Equal(t, _error.BadLogin(), errHTTP)
		assert.Nil(t, userResp)
		assert.Nil(t, profileResp)
	})

	t.Run("ERROR_empty-profile", func(t *testing.T) {
		req := &dto.LoginReq{
			EmailOrUsername: "iban.rama",
			Password:        "rama123",
			RememberMe:      true,
			Oauth2:          false,
		}
		err := userRepo.OpenConn(context.TODO())
		userRepo.OpenConnReturns(nil)
		assert.NoError(t, err)
		defer userRepo.CloseConn()

		_, _ = userRepo.GetUserByEmailOrUsername(context.TODO(), req.EmailOrUsername)
		assert.NoError(t, err)
		userRepo.GetUserByEmailOrUsernameReturns(user, nil)

		_, err = accountApi.GetProfileByUserID(user.ID)
		assert.NoError(t, err)
		accountApi.GetProfileByUserIDReturns(nil, errors.New("BAD GATEWAY"))

		userResp, profileResp, err := authUsecase.Login(context.TODO(), req)
		assert.Error(t, err)
		var errHTTP *model.ErrResponseHTTP
		assert.Equal(t, true, errors.As(err, &errHTTP))
		assert.Equal(t, 502, errHTTP.Code)
		assert.Nil(t, userResp)
		assert.Nil(t, profileResp)
	})
}

func TestAuthUsecaseRegiser(t *testing.T) {
	infra.EnvInit()
	userRepo := &mocks.FakeUserSqlRepo{}
	accessRepo := &mocks.FakeAccessSqlRepo{}

	accountApi := &mocks.FakeAccountApiRepo{}

	authUsecase := _usecase.NewAuthUsecaseImpl(userRepo, accessRepo, accountApi)

	req := &dto.RegisterReq{
		FullName:        "ibanrama",
		Username:        "rama",
		Email:           "ibanrama29@gmail.com",
		Password:        "rama123",
		RePassword:      "rama123",
		EmailVerifiedAt: false,
		AppID:           "1234567",
		Role:            0,
	}

	profile := &model.Profile{
		Code:      201,
		ProfileID: "profileID_1",
		Quote:     "null",
		Profesi:   "null",
	}

	t.Run("SUCCESS", func(t *testing.T) {
		err := userRepo.OpenConn(context.TODO())
		userRepo.OpenConnReturns(nil)
		assert.NoError(t, err)
		defer userRepo.CloseConn()

		err = userRepo.StartTx(context.TODO(), &sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
			ReadOnly:  false,
		})
		assert.NoError(t, err)
		userRepo.StartTxReturns(nil)
		defer func() {
			errEndTx := userRepo.EndTx(err)
			assert.NoError(t, errEndTx)
		}()

		_, _ = userRepo.CheckUserByEmail(context.TODO(), req.Email)
		userRepo.CheckUserByEmailReturns(false, nil)

		_, _ = userRepo.CheckUserByUsername(context.TODO(), req.Username)
		userRepo.CheckUserByUsernameReturns(false, nil)

		id := uuid.NewV4().String()
		userConv, accessConv := converter.RegisterReqToModel(req, id)
		_ = userRepo.CreateUser(context.TODO(), userConv)
		userRepo.CreateUserReturns(nil)
		_, _ = accessRepo.CreateAccess(context.TODO(), accessConv)
		accessRepo.CreateAccessReturns(accessConv, nil)

		profileReq := dto.ProfileReq{
			UserID: userConv.ID,
		}

		profileJson, _ := json.Marshal(profileReq)
		_, _ = accountApi.CreateProfile(profileJson)
		accountApi.CreateProfileReturns(profile, nil)

		userResp, err := authUsecase.Register(context.TODO(), req)
		assert.NoError(t, err)
		assert.NotNil(t, userResp)
		assert.Equal(t, userResp.Email, req.Email)
	})

	t.Run("ERROR_exist-email", func(t *testing.T) {
		err := userRepo.OpenConn(context.TODO())
		userRepo.OpenConnReturns(nil)
		assert.NoError(t, err)
		defer userRepo.CloseConn()

		err = userRepo.StartTx(context.TODO(), &sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
			ReadOnly:  false,
		})
		assert.NoError(t, err)
		userRepo.StartTxReturns(nil)
		defer func() {
			errEndTx := userRepo.EndTx(err)
			assert.NoError(t, errEndTx)
		}()

		_, _ = userRepo.CheckUserByEmail(context.TODO(), req.Email)
		userRepo.CheckUserByEmailReturns(true, nil)

		_, _ = userRepo.CheckUserByUsername(context.TODO(), req.Username)
		userRepo.CheckUserByUsernameReturns(false, nil)

		id := uuid.NewV4().String()
		userConv, accessConv := converter.RegisterReqToModel(req, id)
		_ = userRepo.CreateUser(context.TODO(), userConv)
		userRepo.CreateUserReturns(nil)
		_, _ = accessRepo.CreateAccess(context.TODO(), accessConv)
		accessRepo.CreateAccessReturns(accessConv, nil)

		profileReq := dto.ProfileReq{
			UserID: userConv.ID,
		}

		profileJson, _ := json.Marshal(profileReq)
		_, _ = accountApi.CreateProfile(profileJson)
		accountApi.CreateProfileReturns(profile, nil)

		userResp, err := authUsecase.Register(context.TODO(), req)
		assert.Error(t, err)
		assert.Nil(t, userResp)
		var errHTTP *model.ErrResponseHTTP
		assert.Equal(t, true, errors.As(err, &errHTTP))
		assert.Equal(t, _error.BadExistField("email", "email has been registered"), errHTTP)
	})

	t.Run("ERROR_exist-username", func(t *testing.T) {
		err := userRepo.OpenConn(context.TODO())
		userRepo.OpenConnReturns(nil)
		assert.NoError(t, err)
		defer userRepo.CloseConn()

		err = userRepo.StartTx(context.TODO(), &sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
			ReadOnly:  false,
		})
		assert.NoError(t, err)
		userRepo.StartTxReturns(nil)
		defer func() {
			errEndTx := userRepo.EndTx(err)
			assert.NoError(t, errEndTx)
		}()

		_, _ = userRepo.CheckUserByEmail(context.TODO(), req.Email)
		userRepo.CheckUserByEmailReturns(false, nil)

		_, _ = userRepo.CheckUserByUsername(context.TODO(), req.Username)
		userRepo.CheckUserByUsernameReturns(true, nil)

		id := uuid.NewV4().String()
		userConv, accessConv := converter.RegisterReqToModel(req, id)
		_ = userRepo.CreateUser(context.TODO(), userConv)
		userRepo.CreateUserReturns(nil)
		_, _ = accessRepo.CreateAccess(context.TODO(), accessConv)
		accessRepo.CreateAccessReturns(accessConv, nil)

		profileReq := dto.ProfileReq{
			UserID: userConv.ID,
		}

		profileJson, _ := json.Marshal(profileReq)
		_, _ = accountApi.CreateProfile(profileJson)
		accountApi.CreateProfileReturns(profile, nil)

		userResp, err := authUsecase.Register(context.TODO(), req)
		assert.Error(t, err)
		assert.Nil(t, userResp)
		var errHTTP *model.ErrResponseHTTP
		assert.Equal(t, true, errors.As(err, &errHTTP))
		assert.Equal(t, _error.BadExistField("username", "username has been registered"), errHTTP)
	})

	t.Run("ERROR_bad-account-api", func(t *testing.T) {
		err := userRepo.OpenConn(context.TODO())
		userRepo.OpenConnReturns(nil)
		assert.NoError(t, err)
		defer userRepo.CloseConn()

		err = userRepo.StartTx(context.TODO(), &sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
			ReadOnly:  false,
		})
		assert.NoError(t, err)
		userRepo.StartTxReturns(nil)
		defer func() {
			errEndTx := userRepo.EndTx(err)
			assert.NoError(t, errEndTx)
		}()

		_, _ = userRepo.CheckUserByEmail(context.TODO(), req.Email)
		userRepo.CheckUserByEmailReturns(false, nil)

		_, _ = userRepo.CheckUserByUsername(context.TODO(), req.Username)
		userRepo.CheckUserByUsernameReturns(false, nil)

		id := uuid.NewV4().String()
		userConv, accessConv := converter.RegisterReqToModel(req, id)
		_ = userRepo.CreateUser(context.TODO(), userConv)
		userRepo.CreateUserReturns(nil)
		_, _ = accessRepo.CreateAccess(context.TODO(), accessConv)
		accessRepo.CreateAccessReturns(accessConv, nil)

		profileReq := dto.ProfileReq{
			UserID: userConv.ID,
		}

		profileJson, _ := json.Marshal(profileReq)
		_, _ = accountApi.CreateProfile(profileJson)
		accountApi.CreateProfileReturns(nil, errors.New("bad gateway"))

		userResp, err := authUsecase.Register(context.TODO(), req)
		assert.Error(t, err)
		assert.Nil(t, userResp)
		var errHTTP *model.ErrResponseHTTP
		assert.Equal(t, true, errors.As(err, &errHTTP))
		assert.Equal(t, 502, errHTTP.Code)
	})
}
