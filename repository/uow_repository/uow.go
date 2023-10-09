package uow_repository

import (
	"database/sql"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository"
)

type UnitOfWorkRepositoryImpl struct {
	tx *sql.Tx
	db *sql.DB
}

func NewUnitOfWorkRepositoryImpl(db *sql.DB) repository.UnitOfWorkRepository {
	return &UnitOfWorkRepositoryImpl{
		db: db,
	}
}
