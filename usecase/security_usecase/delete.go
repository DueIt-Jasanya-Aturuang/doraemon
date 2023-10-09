package security_usecase

import (
	"context"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase"
)

func (s *SecurityUsecaseImpl) deletedAllToken(ctx context.Context, userID string) error {
	err := s.securityRepo.StartTx(ctx, repository.LevelReadCommitted(), func() error {
		err := s.securityRepo.DeleteAllByUserID(ctx, userID)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return usecase.InvalidToken
}

func (s *SecurityUsecaseImpl) deletedToken(ctx context.Context, id int, userID string) error {
	err := s.securityRepo.StartTx(ctx, repository.LevelReadCommitted(), func() error {
		err := s.securityRepo.Delete(ctx, id, userID)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
