package _usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/usecase"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra/config"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/converter"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/helper"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util/error"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util/msg"
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
	// open connection dari user repo
	// defer untuk close connection
	err = s.userRepo.OpenConn(ctx)
	if err != nil {
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer s.userRepo.CloseConn()

	// get user by id
	// jika user nya tidak tersedia maka akan return 401
	user, err := s.userRepo.GetUserByID(ctx, req.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn().Msgf("request header user id tidak terdaftar di db")
			return _error.ErrStringDefault(http.StatusUnauthorized)
		}
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}

	// check apakah oldpassword user dan password user di db sama atau gak
	checkPassword := helper.BcryptPasswordCompare(req.OldPassword, user.Password)
	if !checkPassword {
		log.Warn().Msgf("password lama user salah")
		return _error.Err400(map[string][]string{
			"old_password": {
				"password lama anda salah",
			},
		})
	}

	// hasing password request menggunakan bcrypt
	passwordHash, err := helper.BcryptPasswordHash(req.Password)
	if err != nil {
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}

	// start transaction untuk melakukan insert
	// defer untuk commit atau rollback, jika terjadi kesalahan maka akan return 500
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

	// convert request kedalam user model
	// lalu update password user jika error akan return 500
	userConv := converter.ResetPasswordReqToModel(passwordHash, user.ID)
	err = s.userRepo.UpdatePasswordUser(ctx, userConv)
	if err != nil {
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}

	return nil
}

func (s *UserUsecaseImpl) ForgottenPassword(ctx context.Context, req *dto.ForgottenPasswordReq) (url string, err error) {
	// open connection dari user repo
	// defer untuk close connection
	err = s.userRepo.OpenConn(ctx)
	if err != nil {
		return "", _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer s.userRepo.CloseConn()

	// get user by email
	// jika error sql no rows maka akan return 404 bahwa user tidak di temukan
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn().Msgf("user tidak di temukan di database")
			return "", _error.ErrString("USER TIDAK DI TEMUKAN", http.StatusNotFound)
		}
		return "", _error.ErrStringDefault(http.StatusInternalServerError)
	}

	// generate token untuk pergantian di form password, ini akan expired dalam 5 menit
	var jwtModel *model.Jwt
	jwtModelFPT := jwtModel.ForgotPasswordTokenDefault(user.ID)
	forgotPasswordToken, err := helper.GenerateJwtHS256(jwtModelFPT)
	if err != nil {
		return "", _error.ErrStringDefault(http.StatusInternalServerError)
	}

	// set di redis token nya, isinya userid dan exp nya 5 menit
	err = s.redis.Client.Set(ctx, "forgot-password-link:"+req.Email, user.ID, jwtModelFPT.Exp).Err()
	if err != nil {
		log.Err(err).Msg("failed set data in redis")
		return "", _error.ErrStringDefault(http.StatusInternalServerError)
	}

	// akan return link yang berguna untuk melakukan reset password menggunakan method post
	// example : http://example.com/forgot-password?email=ibanrama29@gmail.com&token=tokenrama
	link := fmt.Sprintf("%s/forgot-password?email=%s&token=%s", config.AppAuthApi, req.Email, forgotPasswordToken)
	return link, nil
}

func (s *UserUsecaseImpl) ResetForgottenPassword(ctx context.Context, req *dto.ResetForgottenPasswordReq) (err error) {
	// claim token user untuk melakukan reset password
	// jika error makan akan return 401
	claims, err := helper.ClaimsJwtHS256(req.Token, config.DefaultKey)
	if err != nil {
		log.Err(err).Msg("failed claim token user")
		return _error.ErrString("INVALID YOUR TOKEN", http.StatusUnauthorized)
	}

	// type assertion sub token, jika error akan return 401
	userID, ok := claims["sub"].(string)
	if !ok {
		log.Warn().Msgf("gagal type assertion pada sub token")
		return _error.ErrString("INVALID YOUR TOKEN", http.StatusUnauthorized)
	}

	// get userid yang sudah di redis
	getUserID, err := s.redis.Client.Get(ctx, "forgot-password-link:"+req.Email).Result()
	if err != nil {
		log.Err(err).Msg("user mencoba untuk get data di redis")
		return _error.ErrStringDefault(http.StatusUnauthorized)
	}

	// jika userid di sub token dengan userid di redis tidak match maka akan return 401
	if userID != getUserID {
		log.Warn().Msgf("userid di redis sama userid di token tidak sama")
		return _error.ErrString("INVALID YOUR TOKEN", http.StatusUnauthorized)
	}

	// hash bcrypt password yang baru
	passwordHash, err := helper.BcryptPasswordHash(req.Password)
	if err != nil {
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}

	// open connection dari user repo
	// defer untuk close connection
	err = s.userRepo.OpenConn(ctx)
	if err != nil {
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer s.userRepo.CloseConn()

	// start transaction
	// defer untuk commit atau rollback, jika error akan return 500
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

	// convert request to user model untuk update user
	// dan update password user berdasarkan user convert tdi
	userConv := converter.ResetPasswordReqToModel(passwordHash, userID)
	err = s.userRepo.UpdatePasswordUser(ctx, userConv)
	if err != nil {
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}

	err = s.redis.Client.Del(ctx, "forgot-password-link:"+req.Email).Err()
	if err != nil {
		log.Err(err).Msg(_msg.LogErrDelRedisClient)
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}
	return nil
}

func (s *UserUsecaseImpl) ActivasiAccount(ctx context.Context, email string) (resp *dto.ActivasiAccountResp, err error) {
	// open connection dari user repo
	// defer untuk close connection
	err = s.userRepo.OpenConn(ctx)
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer s.userRepo.CloseConn()

	// get user by email, jika email (ErrNoRows) tidak tersedia maka akan return 404
	// jika error maka akan return 500
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn().Msgf("ketika user mau aktivasi email nya tidak ditemukan di database")
			return nil, _error.ErrString("USER TIDAK DITEMUKAN", 404)
		}
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	// check apakah user sudah aktiv apa belum
	// kalau sudah maka akan return 400
	if user.EmailVerifiedAt {
		log.Warn().Msgf("user sudah melakukan activasi akun tapi mencoba aktivasi kembali")
		return nil, _error.Err400(map[string][]string{
			"email": {
				"permintaan anda tidak dapat di proses, email anda sudah di aktivasi silahkan login",
			},
		})
	}

	// start transaction
	// defer untuk commit atau rollback, jika error akan return 500
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

	// usermodel untuk memperbarui data user kalau dia sudah di aktivasi
	// lalu insert update user, jika error kan return 500
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

	// response untuk fe memperbarui bahwa user sudah di aktivasi
	resp = &dto.ActivasiAccountResp{
		EmailVerifiedAt: true,
	}

	return resp, nil
}
