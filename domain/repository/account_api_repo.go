package repository

import (
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
)

//counterfeiter:generate -o ./../mocks . AccountApiRepo
type AccountApiRepo interface {
	CreateProfile(data []byte) (*model.Profile, error)
	GetProfileByUserID(userID string) (*model.Profile, error)
}
