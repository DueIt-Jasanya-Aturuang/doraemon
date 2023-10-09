package uow_repository

import (
	"database/sql"
	"fmt"

	"github.com/rs/zerolog/log"
)

func (u *UnitOfWorkRepositoryImpl) GetDB() (*sql.DB, error) {
	if u.db == nil {
		err := fmt.Errorf("no Connection Database Available")
		log.Warn().Msg(err.Error())
		return nil, err
	}

	return u.db, nil
}
