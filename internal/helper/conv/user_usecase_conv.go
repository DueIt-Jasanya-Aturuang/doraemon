package conv

import (
	"database/sql"
	"time"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
)

func ResetPasswordReqToModel(password string, userID string) *model.User {
	return &model.User{
		ID:        userID,
		Password:  password,
		UpdatedAt: time.Now().Unix(),
		UpdatedBy: sql.NullString{String: userID, Valid: true},
	}
}
