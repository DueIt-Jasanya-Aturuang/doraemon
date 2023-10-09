package user_usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase"
)

func (u *UserUsecaseImpl) ChangePassword(ctx context.Context, req *usecase.RequestChangePassword) error {
	repository.GetUserByID = req.UserID
	user, err := u.userRepo.Get(ctx, repository.GetUserByID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return usecase.InvalidUserID
		}
		return err
	}

	checkPass := usecase.BcryptPasswordCompare(req.OldPassword, user.Password)
	if !checkPass {
		log.Warn().Msgf("password lama user salah")
		return usecase.InvalidOldPassword
	}

	newPass, err := usecase.BcryptPasswordHash(req.Password)
	if err != nil {
		return err
	}

	err = u.userRepo.StartTx(ctx, repository.LevelReadCommitted(), func() error {
		userConv := usecase.ChangePasswordRequestToModel(newPass, user.ID)
		err = u.userRepo.UpdatePassword(ctx, userConv)
		return err
	})

	return err
}

func (u *UserUsecaseImpl) ChangeUsername(ctx context.Context, req *usecase.RequestChangeUsername) error {
	repository.GetUserByID = req.UserID
	user, err := u.userRepo.Get(ctx, repository.GetUserByID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return usecase.InvalidUserID
		}
		return err
	}

	if user.Username == req.Username {
		return nil
	}

	repository.CheckUserByUsername = req.Username
	exist, err := u.userRepo.Check(ctx, repository.CheckUserByUsername)
	if err != nil {
		return err
	}
	if exist {
		return usecase.UsernameIsExist
	}

	err = u.userRepo.StartTx(ctx, repository.LevelReadCommitted(), func() error {
		userConv := usecase.ChangeUsernameRequestToModel(req.Username, user.ID)
		err = u.userRepo.UpdateUsername(ctx, userConv)
		return err
	})

	return err
}

func (u *UserUsecaseImpl) ActivasiAccount(ctx context.Context, email string) (bool, error) {
	repository.GetUserByEmail = email
	user, err := u.userRepo.Get(ctx, repository.GetUserByEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, usecase.InvalidUserID
		}
		return false, err
	}

	if user.EmailVerifiedAt {
		log.Warn().Msgf("user sudah melakukan activasi akun tapi mencoba aktivasi kembali")
		return false, usecase.EmailIsActivited
	}

	err = u.userRepo.StartTx(ctx, repository.LevelReadCommitted(), func() error {
		userConv := usecase.ActivasiAccountRequestToModel(user.ID)
		err = u.userRepo.UpdateActivasi(ctx, userConv)
		return err
	})

	if err != nil {
		return false, err
	}

	return true, nil
}
