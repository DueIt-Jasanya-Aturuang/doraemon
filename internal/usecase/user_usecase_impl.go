package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infrastructures/config"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/helper"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/helper/conv"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/encryption"
	_error "github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/error"
)

type UserUsecaseImpl struct {
	userRepo repository.UserSqlRepo
	redis    *config.RedisImpl
}

func NewUserUsecaseImpl(
	userRepo repository.UserSqlRepo,
	redis *config.RedisImpl,
) usecase.UserUsecase {
	return &UserUsecaseImpl{
		userRepo: userRepo,
		redis:    redis,
	}
}

func (s *UserUsecaseImpl) ResetPassword(ctx context.Context, req *dto.ResetPasswordReq) (err error) {
	err = s.userRepo.OpenConn(ctx)
	if err != nil {
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer s.userRepo.CloseConn()

	user, err := s.userRepo.GetUserByID(ctx, req.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return _error.ErrStringDefault(http.StatusUnauthorized)
		}
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}

	checkPassword := helper.BcryptPasswordCompare(req.OldPassword, user.Password)
	if !checkPassword {
		return _error.Err400(map[string][]string{
			"old_password": {
				"password lama anda salah",
			},
		})
	}

	passwordHash, err := helper.BcryptPasswordHash(req.Password)
	if err != nil {
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}

	err = s.userRepo.StartTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer func() {
		errEndTx := s.userRepo.EndTx(err)
		if errEndTx != nil {
			err = _error.ErrStringDefault(http.StatusInternalServerError)
		}
	}()

	userConv := conv.ResetPasswordReqToModel(passwordHash, user.ID)
	err = s.userRepo.UpdatePasswordUser(ctx, userConv)
	if err != nil {
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}

	return nil
}

func (s *UserUsecaseImpl) ForgottenPassword(ctx context.Context, req *dto.ForgottenPasswordReq) (url string, err error) {
	err = s.userRepo.OpenConn(ctx)
	if err != nil {
		return "", _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer s.userRepo.CloseConn()

	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", _error.ErrString("USER TIDAK DI TEMUKAN", http.StatusNotFound)
		}
		return "", _error.ErrStringDefault(http.StatusInternalServerError)
	}

	tokenID := uuid.NewV4().String()
	var jwtModel *model.Jwt

	jwtModelFPT := jwtModel.ForgotPasswordTokenDefault(tokenID, user.ID)
	forgotPasswordToken, err := helper.GenerateJwtHS256(jwtModelFPT)
	if err != nil {
		return "", _error.ErrStringDefault(http.StatusInternalServerError)
	}

	_, err = s.redis.Client.Set(ctx, "forgot-password-link:"+req.Email, tokenID, jwtModelFPT.Exp).Result()
	if err != nil {
		return "", _error.ErrStringDefault(http.StatusInternalServerError)
	}

	link := fmt.Sprintf("%s/forgot-password?email=%s&token=%s", config.AppAuthApi, req.Email, forgotPasswordToken)
	return link, nil
}

func (s *UserUsecaseImpl) ResetForgottenPassword(ctx context.Context, req *dto.ResetForgottenPasswordReq) (err error) {
	claims, err := helper.ClaimsJwtHS256(req.Token, config.DefaultKey)
	if err != nil {
		return _error.ErrString("INVALID YOUR TOKEN", http.StatusUnauthorized)
	}

	sub, err := encryption.DecryptStringCFB(claims["sub"].(string), config.AesCFB)
	if err != nil {
		return nil
	}

	tokenArray := strings.Split(sub, ":")
	if len(tokenArray) != 3 {
		return _error.ErrString("INVALID YOUR TOKEN", http.StatusUnauthorized)
	}

	if tokenArray[2] != "forgot-password" {
		return _error.ErrString("INVALID YOUR TOKEN", http.StatusUnauthorized)
	}

	tokenID := tokenArray[0]

	getTokenID, err := s.redis.Client.Get(ctx, "forgot-password-link:"+req.Email).Result()
	if err != nil {
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}

	if tokenID != getTokenID {
		return _error.ErrString("INVALID YOUR TOKEN", http.StatusUnauthorized)
	}

	passwordHash, err := helper.BcryptPasswordHash(req.Password)
	if err != nil {
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}

	err = s.userRepo.OpenConn(ctx)
	if err != nil {
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer s.userRepo.CloseConn()

	err = s.userRepo.StartTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer func() {
		errEndTx := s.userRepo.EndTx(err)
		if errEndTx != nil {
			err = _error.ErrStringDefault(http.StatusInternalServerError)
		}
	}()

	userID := tokenArray[1]
	userConv := conv.ResetPasswordReqToModel(passwordHash, userID)
	
	err = s.userRepo.UpdatePasswordUser(ctx, userConv)
	if err != nil {
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}

	return nil
}

func (s *UserUsecaseImpl) ActivasiAccount(ctx context.Context, email string) (resp *dto.ActivasiAccountResp, err error) {
	err = s.userRepo.OpenConn(ctx)
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer s.userRepo.CloseConn()

	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, _error.ErrString("USER TIDAK DITEMUKAN", 404)
		}
	}

	if user.EmailVerifiedAt {
		return nil, _error.Err400(map[string][]string{
			"email": {
				"permintaan anda tidak dapat di proses, email anda sudah di aktivasi silahkan login",
			},
		})
	}

	err = s.userRepo.StartTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	defer func() {
		errEndTx := s.userRepo.EndTx(err)
		if errEndTx != nil {
			err = _error.ErrStringDefault(http.StatusInternalServerError)
		}
	}()

	userModel := &model.User{
		ID:              user.ID,
		EmailVerifiedAt: true,
		UpdatedAt:       time.Now().Unix(),
		UpdatedBy:       sql.NullString{String: user.ID, Valid: true},
	}

	err = s.userRepo.UpdateActivasiUser(ctx, userModel)
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	resp = &dto.ActivasiAccountResp{
		EmailVerifiedAt: true,
	}

	return resp, nil
}
