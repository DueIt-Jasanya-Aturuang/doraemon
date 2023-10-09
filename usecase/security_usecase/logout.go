package security_usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase"
)

func (s *SecurityUsecaseImpl) Logout(ctx context.Context, req *usecase.RequestLogout) error {
	token, err := s.securityRepo.GetByAccessToken(ctx, req.Token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn().Msg("user mencoba untuk logout menggunakan token yang lama")
			return nil
		}
		return err
	}

	err = s.deletedToken(ctx, token.ID, token.UserID)
	if err != nil {
		return err
	}

	return nil
}
