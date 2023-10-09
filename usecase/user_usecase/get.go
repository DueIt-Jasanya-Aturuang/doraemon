package user_usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase"
)

func (u *UserUsecaseImpl) GetUserByID(ctx context.Context, id string) (*usecase.ResponseUser, error) {
	repository.GetUserByID = id
	user, err := u.userRepo.Get(ctx, repository.GetUserByID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, usecase.InvalidUserID
		}
		return nil, err
	}

	resp := usecase.UserModelToResponse(user)

	return resp, nil
}
