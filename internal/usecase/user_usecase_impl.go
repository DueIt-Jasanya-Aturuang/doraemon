package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infrastructures/config"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/helper"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/helper/conv"
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
			log.Warn().Msgf("request header user id tidak terdaftar di db")
			return _error.ErrStringDefault(http.StatusUnauthorized)
		}
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}

	checkPassword := helper.BcryptPasswordCompare(req.OldPassword, user.Password)
	if !checkPassword {
		log.Warn().Msgf("password lama user salah")
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
			log.Warn().Msgf("user tidak di temukan di database")
			return "", _error.ErrString("USER TIDAK DI TEMUKAN", http.StatusNotFound)
		}
		return "", _error.ErrStringDefault(http.StatusInternalServerError)
	}

	var jwtModel *model.Jwt

	jwtModelFPT := jwtModel.ForgotPasswordTokenDefault(user.ID)
	forgotPasswordToken, err := helper.GenerateJwtHS256(jwtModelFPT)
	if err != nil {
		return "", _error.ErrStringDefault(http.StatusInternalServerError)
	}

	err = s.redis.Client.Set(ctx, "forgot-password-link:"+req.Email, user.ID, jwtModelFPT.Exp).Err()
	if err != nil {
		log.Err(err).Msg("failed set data in redis")
		return "", _error.ErrStringDefault(http.StatusInternalServerError)
	}

	link := fmt.Sprintf("%s/forgot-password?email=%s&token=%s", config.AppAuthApi, req.Email, forgotPasswordToken)
	return link, nil
}

func (s *UserUsecaseImpl) ResetForgottenPassword(ctx context.Context, req *dto.ResetForgottenPasswordReq) (err error) {
	claims, err := helper.ClaimsJwtHS256(req.Token, config.DefaultKey)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			log.Err(err).Msg("token user ketika di claim suda expired")
		}
		log.Err(err).Msg("failed claim token user")
		return _error.ErrString("INVALID YOUR TOKEN", http.StatusUnauthorized)
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		log.Warn().Msgf("gagal type assertion pada sub token")
		return _error.ErrString("INVALID YOUR TOKEN", http.StatusUnauthorized)
	}

	getUserID, err := s.redis.Client.Get(ctx, "forgot-password-link:"+req.Email).Result()
	if err != nil {
		log.Err(err).Msg("failed get data in redis")
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}

	if userID != getUserID {
		log.Warn().Msgf("userid di redis sama userid di token tidak sama")
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
			log.Warn().Msgf("ketika user mau aktivasi email nya tidak ditemukan di database")
			return nil, _error.ErrString("USER TIDAK DITEMUKAN", 404)
		}
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	if user.EmailVerifiedAt {
		log.Warn().Msgf("user sudah melakukan activasi akun tapi mencoba aktivasi kembali")
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
