package repository

import (
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
)

type AccountApiRepo interface {
	CreateProfile(data []byte) (*model.Profile, error)
}
