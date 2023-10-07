package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase/converter"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase/helper"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

type UserUsecaseImpl struct {
	userRepo domain.UserRepository
	redis    *infra.RedisImpl
}

func NewUserUsecaseImpl(
	userRepo domain.UserRepository,
	redis *infra.RedisImpl,
) domain.UserUsecase {
	return &UserUsecaseImpl{
		userRepo: userRepo,
		redis:    redis,
	}
}

func (u *UserUsecaseImpl) ChangePassword(ctx context.Context, req *domain.RequestChangePassword) error {
	domain.GetUserByID = req.UserID
	user, err := u.userRepo.Get(ctx, domain.GetUserByID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return InvalidUserID
		}
		return err
	}

	checkPass := helper.BcryptPasswordCompare(req.OldPassword, user.Password)
	if !checkPass {
		log.Warn().Msgf("password lama user salah")
		return InvalidOldPassword
	}

	newPass, err := helper.BcryptPasswordHash(req.Password)
	if err != nil {
		return err
	}

	err = u.userRepo.StartTx(ctx, helper.LevelReadCommitted(), func() error {
		userConv := converter.ChangePasswordReqToModel(newPass, user.ID)
		err = u.userRepo.UpdatePassword(ctx, userConv)
		return err
	})

	return err
}

func (u *UserUsecaseImpl) ChangeUsername(ctx context.Context, req *domain.RequestChangeUsername) error {
	domain.GetUserByID = req.UserID
	user, err := u.userRepo.Get(ctx, domain.GetUserByID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return InvalidUserID
		}
		return err
	}

	if user.Username == req.Username {
		return nil
	}

	domain.CheckUserByUsername = req.Username
	exist, err := u.userRepo.Check(ctx, domain.CheckUserByUsername)
	if err != nil {
		return err
	}
	if exist {
		return UsernameIsExist
	}

	err = u.userRepo.StartTx(ctx, helper.LevelReadCommitted(), func() error {
		userConv := converter.ChangeUsernameReqToModel(req.Username, user.ID)
		err = u.userRepo.UpdateUsername(ctx, userConv)
		return err
	})

	return err
}

func (u *UserUsecaseImpl) ForgottenPassword(ctx context.Context, req *domain.RequestForgottenPassword) (string, error) {
	domain.GetUserByEmail = req.Email
	user, err := u.userRepo.Get(ctx, domain.GetUserByEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", InvalidUserID
		}
		return "", err
	}

	var jwtModel *helper.Jwt
	fpTokenModel := jwtModel.ForgotPasswordTokenDefault(user.ID)

	fpToken, err := helper.GenerateJwtHS256(fpTokenModel)
	if err != nil {
		return "", err
	}

	err = u.redis.Client.Set(ctx, util.ForgotPasswordLink+":"+req.Email, user.ID, fpTokenModel.Exp).Err()
	if err != nil {
		log.Warn().Msgf(util.LogErrRedisClientSet, err)
		return "", err
	}

	link := fmt.Sprintf("%s/forgot-password?email=%s&token=%s", infra.AppAuthApi, req.Email, fpToken)
	return link, nil
}

func (u *UserUsecaseImpl) ResetForgottenPassword(ctx context.Context, req *domain.RequestResetForgottenPassword) error {
	claims, err := helper.ClaimsJwtHS256(req.Token, infra.DefaultKey)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return TokenExpired
		}
		log.Warn().Msgf("failed claims token | err : %v", err)
		return InvalidToken
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		log.Warn().Msgf("cannot assetion claims sub to string | data : %v", claims["sub"])
		return InvalidToken
	}

	getUserID, err := u.redis.Client.Get(ctx, util.ForgotPasswordLink+":"+req.Email).Result()
	if err != nil {
		log.Warn().Msgf(util.LogErrRedisClientGet, err)
		return TokenExpired
	}

	if userID != getUserID {
		log.Warn().Msgf("userid di redis sama userid di token tidak sama")
		return InvalidToken
	}

	newPass, err := helper.BcryptPasswordHash(req.Password)
	if err != nil {
		return err
	}

	err = u.userRepo.StartTx(ctx, helper.LevelReadCommitted(), func() error {
		userConv := converter.ChangePasswordReqToModel(newPass, userID)
		err = u.userRepo.UpdatePassword(ctx, userConv)
		return err
	})

	return err
}

func (u *UserUsecaseImpl) ActivasiAccount(ctx context.Context, email string) (*domain.ResponseActivasiAccount, error) {
	domain.GetUserByEmail = email
	user, err := u.userRepo.Get(ctx, domain.GetUserByEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, InvalidUserID
		}
		return nil, err
	}

	if user.EmailVerifiedAt {
		log.Warn().Msgf("user sudah melakukan activasi akun tapi mencoba aktivasi kembali")
		return nil, EmailIsActivited
	}

	err = u.userRepo.StartTx(ctx, helper.LevelReadCommitted(), func() error {
		userConv := converter.ActivasiAccountReqToModel(user.ID)
		err = u.userRepo.UpdateActivasi(ctx, userConv)
		return err
	})

	if err != nil {
		return nil, err
	}

	resp := &domain.ResponseActivasiAccount{
		EmailVerifiedAt: true,
	}
	return resp, nil
}

func (u *UserUsecaseImpl) GetUserByID(ctx context.Context, id string) (*domain.ResponseUser, error) {
	domain.GetUserByID = id
	user, err := u.userRepo.Get(ctx, domain.GetUserByID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, InvalidUserID
		}
		return nil, err
	}

	resp := converter.UserModelToResp(user)

	return resp, nil
}
