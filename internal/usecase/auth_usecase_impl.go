package usecase

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infrastructures/config"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/helper"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/helper/conv"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/helper/validation"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/encryption"
	_error "github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/error"
)

type AuthUsecaseImpl struct {
	userRepo     repository.UserSqlRepo
	appRepo      repository.AppSqlRepo
	accessRepo   repository.AccessSqlRepo
	securityRepo repository.SecuritySqlRepo
	timeout      time.Duration
}

func NewAuthUsecaseImpl(
	userRepo repository.UserSqlRepo,
	appRepo repository.AppSqlRepo,
	accessRepo repository.AccessSqlRepo,
	securityRepo repository.SecuritySqlRepo,
	timeout time.Duration,
) usecase.AuthUsecase {
	return &AuthUsecaseImpl{
		userRepo:     userRepo,
		appRepo:      appRepo,
		accessRepo:   accessRepo,
		securityRepo: securityRepo,
		timeout:      timeout,
	}
}

func (a *AuthUsecaseImpl) Login(c context.Context, req *dto.LoginReq) (userResp *dto.UserResp, err error) {
	ctx, cancel := context.WithTimeout(c, a.timeout)
	defer cancel()

	err = validation.LoginValidation(req)
	if err != nil {
		return nil, err
	}

	err = a.userRepo.OpenConn(ctx)
	if err != nil {
		return nil, _error.ErrString("INTERNAL SERVER ERROR", 500)
	}
	defer a.userRepo.CloseConn()

	exists, err := a.appRepo.CheckAppByID(ctx, req.AppId)
	if err != nil {
		return nil, _error.ErrString("INTERNAL SERVER ERROR", 500)
	}
	if !exists {
		log.Warn().Msgf("app id is not registered: %s", req.AppId)
		return nil, _error.ErrString("FORBIDDEN", 403)
	}

	user, err := a.userRepo.GetUserByEmailOrUsername(ctx, req.EmailOrUsername)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, _error.BadLogin()
		}
		return nil, _error.ErrString("INTERNAL SERVER ERROR", 500)
	}

	if !req.Oauth2 {
		checkPassword := helper.BcryptPasswordCompare(req.Password, user.Password)
		if !checkPassword {
			return nil, _error.BadLogin()
		}
	}

	tokenID := uuid.NewV4().String()
	var jwtModel *model.Jwt

	jwtModelAT := jwtModel.AccessTokenDefault(tokenID, user.ID, req.RememberMe)
	accessToken, err := helper.GenerateJwtHS256(jwtModelAT)
	if err != nil {
		return nil, _error.ErrString("INTERNAL SERVER ERROR", 500)
	}

	jwtModelRT := jwtModel.RefreshTokenDefault(tokenID, user.ID, req.RememberMe)
	refreshToken, err := helper.GenerateJwtHS256(jwtModelRT)
	if err != nil {
		return nil, _error.ErrString("INTERNAL SERVER ERROR", 500)
	}

	err = a.userRepo.StartTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, _error.ErrString("INTERNAL SERVER ERROR", 500)
	}
	defer func() {
		errEndTx := a.userRepo.EndTx(err)
		if errEndTx != nil {
			err = _error.ErrString("INTERNAL SERVER ERROR", 500)
			userResp = nil
		}
	}()

	_, err = a.securityRepo.CreateToken(ctx, &model.Token{
		ID:     tokenID,
		UserID: user.ID,
		AppID:  req.AppId,
		Token:  refreshToken,
	})
	if err != nil {
		return nil, _error.ErrString("INTERNAL SERVER ERROR", 500)
	}

	emailFormat := util.EmailFormat(user.Email)
	userResp = conv.LoginModelToResp(user, accessToken, emailFormat)

	return userResp, nil
}

func (a *AuthUsecaseImpl) Logout(c context.Context, req *dto.LogoutReq) error {
	ctx, cancel := context.WithTimeout(c, a.timeout)
	defer cancel()

	claims, err := helper.ClaimsJwtHS256(config.AccessTokenKeyHS, req.Token)
	if err != nil {
		return nil
	}
	sub, err := encryption.DecryptStringCFB(claims["sub"].(string), config.AesCFB)
	if err != nil {
		return nil
	}

	tokenID := strings.Split(sub, ":")[0]
	userID := strings.Split(sub, ":")[1]

	err = a.userRepo.OpenConn(ctx)
	if err != nil {
		return _error.ErrString("INTERNAL SERVER ERROR", 500)
	}
	defer a.userRepo.CloseConn()

	err = a.userRepo.StartTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return _error.ErrString("INTERNAL SERVER ERROR", 500)
	}
	defer func() {
		errEndTx := a.userRepo.EndTx(err)
		if errEndTx != nil {
			err = _error.ErrString("INTERNAL SERVER ERROR", 500)
		}
	}()

	err = a.securityRepo.DeleteToken(ctx, tokenID, userID)
	if err != nil {
		return _error.ErrString("INTERNAL SERVER ERROR", 500)
	}

	return nil
}

func (a *AuthUsecaseImpl) Register(c context.Context, req *dto.RegisterReq) (userResp *dto.UserResp, err error) {
	ctx, cancel := context.WithTimeout(c, a.timeout)
	defer cancel()

	err = validation.RegisterValidation(req)
	if err != nil {
		return nil, err
	}

	err = a.userRepo.OpenConn(ctx)
	if err != nil {
		return nil, _error.ErrString("INTERNAL SERVER ERROR", 500)
	}
	defer a.userRepo.CloseConn()

	exists, err := a.appRepo.CheckAppByID(ctx, req.AppId)
	if err != nil {
		return nil, _error.ErrString("INTERNAL SERVER ERROR", 500)
	}
	if !exists {
		log.Warn().Msgf("app id is not registered: %s", req.AppId)
		return nil, _error.ErrString("UNAUTHORIZATION", 403)
	}

	exists, err = a.userRepo.CheckUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, _error.ErrString("INTERNAL SERVER ERROR", 500)
	}
	if exists {
		return nil, _error.Err400(map[string][]string{
			"email": {
				"email has been registered",
			},
		})
	}

	exists, err = a.userRepo.CheckUserByUsername(ctx, req.Email)
	if err != nil {
		return nil, _error.ErrString("INTERNAL SERVER ERROR", 500)
	}
	if exists {
		return nil, _error.Err400(map[string][]string{
			"username": {
				"username has been registered",
			},
		})
	}

	passwordHash, err := helper.BcryptPasswordHash(req.Password)
	if err != nil {
		return nil, _error.ErrString("INTERNAL SERVER ERROR", 500)
	}
	req.Password = passwordHash

	err = a.userRepo.StartTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, _error.ErrString("INTERNAL SERVER ERROR", 500)
	}
	defer func() {
		errEndTx := a.userRepo.EndTx(err)
		if errEndTx != nil {
			err = _error.ErrString("INTERNAL SERVER ERROR", 500)
			userResp = nil
		}
	}()

	id := uuid.NewV4().String()
	userConv, accessConv := conv.RegisterReqToModel(req, id)

	user, err := a.userRepo.CreateUser(ctx, userConv)
	if err != nil {
		return nil, _error.ErrString("INTERNAL SERVER ERROR", 500)
	}

	_, err = a.accessRepo.CreateAccess(ctx, accessConv)
	if err != nil {
		return nil, _error.ErrString("INTERNAL SERVER ERROR", 500)
	}

	emailFormat := util.EmailFormat(user.Email)
	userResp = conv.RegisterModelToResp(user, emailFormat)

	return userResp, nil
}
