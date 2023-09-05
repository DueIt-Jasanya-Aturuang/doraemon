package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infrastructures/config"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/helper"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/helper/conv"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/encryption"
	_error "github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/error"
)

type AuthUsecaseImpl struct {
	userRepo     repository.UserSqlRepo
	accessRepo   repository.AccessSqlRepo
	securityRepo repository.SecuritySqlRepo
	accountApi   repository.AccountApiRepo
}

func NewAuthUsecaseImpl(
	userRepo repository.UserSqlRepo,
	accessRepo repository.AccessSqlRepo,
	securityRepo repository.SecuritySqlRepo,
	accountApi repository.AccountApiRepo,
) usecase.AuthUsecase {
	return &AuthUsecaseImpl{
		userRepo:     userRepo,
		accessRepo:   accessRepo,
		securityRepo: securityRepo,
		accountApi:   accountApi,
	}
}

func (a *AuthUsecaseImpl) Login(ctx context.Context, req *dto.LoginReq) (userResp *dto.UserResp, profileResp *dto.ProfileResp, err error) {
	err = a.userRepo.OpenConn(ctx)
	if err != nil {
		return nil, nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer a.userRepo.CloseConn()

	user, err := a.userRepo.GetUserByEmailOrUsername(ctx, req.EmailOrUsername)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, _error.BadLogin()
		}
		return nil, nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	if !req.Oauth2 {
		checkPassword := helper.BcryptPasswordCompare(req.Password, user.Password)
		if !checkPassword {
			return nil, nil, _error.BadLogin()
		}
	}

	profile, err := a.accountApi.GetProfileByUserID(user.ID)
	if err != nil {
		log.Err(err).Msg("error account service")
		return nil, nil, _error.ErrStringDefault(http.StatusBadGateway)
	}

	emailFormat := util.EmailFormat(user.Email)
	userResp, profileResp = conv.LoginModelToResp(user, profile, emailFormat)

	return userResp, profileResp, nil
}

func (a *AuthUsecaseImpl) Logout(ctx context.Context, req *dto.LogoutReq) error {
	claims, err := helper.ClaimsJwtHS256(config.AccessTokenKeyHS, req.Token)
	if err != nil {
		if !errors.Is(err, jwt.ErrTokenExpired) {
			return nil
		}
	}

	sub, err := encryption.DecryptStringCFB(claims["sub"].(string), config.AesCFB)
	if err != nil {
		return nil
	}

	tokenID := strings.Split(sub, ":")[0]
	userID := strings.Split(sub, ":")[1]

	err = a.userRepo.OpenConn(ctx)
	if err != nil {
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer a.userRepo.CloseConn()

	err = a.userRepo.StartTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer func() {
		errEndTx := a.userRepo.EndTx(err)
		if errEndTx != nil {
			err = _error.ErrStringDefault(http.StatusInternalServerError)
		}
	}()

	err = a.securityRepo.DeleteToken(ctx, tokenID, userID)
	if err != nil {
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}

	return nil
}

func (a *AuthUsecaseImpl) Register(ctx context.Context, req *dto.RegisterReq) (userResp *dto.UserResp, profileResp *dto.ProfileResp, err error) {
	err = a.userRepo.OpenConn(ctx)
	if err != nil {
		return nil, nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer a.userRepo.CloseConn()

	exists, err := a.userRepo.CheckUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}
	if exists {
		return nil, nil, _error.BadUsername("username has been registered")
	}

	exists, err = a.userRepo.CheckUserByUsername(ctx, req.Email)
	if err != nil {
		return nil, nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}
	if exists {
		return nil, nil, _error.BadUsername("username has been registered")
	}

	passwordHash, err := helper.BcryptPasswordHash(req.Password)
	if err != nil {
		return nil, nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}
	req.Password = passwordHash

	err = a.userRepo.StartTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer func() {
		errEndTx := a.userRepo.EndTx(err)
		if errEndTx != nil {
			err = _error.ErrStringDefault(http.StatusInternalServerError)
			userResp = nil
			profileResp = nil
		}
	}()

	id := uuid.NewV4().String()
	userConv, accessConv := conv.RegisterReqToModel(req, id)

	err = a.userRepo.CreateUser(ctx, userConv)
	if err != nil {
		return nil, nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	_, err = a.accessRepo.CreateAccess(ctx, accessConv)
	if err != nil {
		return nil, nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	profileReq := dto.ProfileReq{
		UserID: userConv.ID,
	}
	profileJson, err := json.Marshal(profileReq)
	if err != nil {
		log.Err(err).Msg("failed marshal profile req")
		return nil, nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	profile, err := a.accountApi.CreateProfile(profileJson)
	if err != nil {
		log.Err(err).Msg("error account service")
		return nil, nil, _error.ErrStringDefault(http.StatusBadGateway)
	}

	emailFormat := util.EmailFormat(userConv.Email)
	userResp, profileResp = conv.RegisterModelToResp(userConv, profile, emailFormat)

	return userResp, profileResp, nil
}
