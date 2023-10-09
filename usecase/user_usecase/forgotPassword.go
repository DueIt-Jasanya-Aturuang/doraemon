package user_usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase/jwt_usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

func (u *UserUsecaseImpl) ForgottenPassword(ctx context.Context, req *usecase.RequestForgottenPassword) (string, error) {
	repository.GetUserByEmail = req.Email
	user, err := u.userRepo.Get(ctx, repository.GetUserByEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", usecase.InvalidUserID
		}
		return "", err
	}

	var jwtModel *jwt_usecase.Jwt
	fpTokenModel := jwtModel.ForgotPasswordTokenDefault(user.ID)

	fpToken, err := jwt_usecase.GenerateJwtHS256(fpTokenModel)
	if err != nil {
		return "", err
	}

	err = u.redis.Client.Set(ctx, util.ForgotPasswordLink+":"+req.Email, user.ID, fpTokenModel.Exp).Err()
	if err != nil {
		log.Warn().Msgf(util.LogErrRedisClientSet, err)
		return "", err
	}

	link := fmt.Sprintf("%s/forgot-password?email=%s&token=%s", infra.BaseUrlAuthService, req.Email, fpToken)
	return link, nil
}

func (u *UserUsecaseImpl) ResetForgottenPassword(ctx context.Context, req *usecase.RequestResetForgottenPassword) error {
	claims, err := jwt_usecase.ClaimsJwtHS256(req.Token, infra.DefaultKey)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return usecase.TokenExpired
		}
		log.Warn().Msgf("failed claims token | err : %v", err)
		return usecase.InvalidToken
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		log.Warn().Msgf("cannot assetion claims sub to string | data : %v", claims["sub"])
		return usecase.InvalidToken
	}

	getUserID, err := u.redis.Client.Get(ctx, util.ForgotPasswordLink+":"+req.Email).Result()
	if err != nil {
		log.Warn().Msgf(util.LogErrRedisClientGet, err)
		return usecase.TokenExpired
	}

	if userID != getUserID {
		log.Warn().Msgf("userid di redis sama userid di token tidak sama")
		return usecase.InvalidToken
	}

	newPass, err := usecase.BcryptPasswordHash(req.Password)
	if err != nil {
		return err
	}

	err = u.userRepo.StartTx(ctx, repository.LevelReadCommitted(), func() error {
		userConv := usecase.ChangePasswordRequestToModel(newPass, userID)
		err = u.userRepo.UpdatePassword(ctx, userConv)
		return err
	})

	return err
}
