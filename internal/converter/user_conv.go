package converter

import (
	"time"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/helper"
)

func ResetPasswordReqToModel(password string, userID string) *domain.User {
	return &domain.User{
		ID:       userID,
		Password: password,
		AuditInfo: domain.AuditInfo{
			UpdatedAt: time.Now().Unix(),
			UpdatedBy: helper.NewNullString(userID),
		},
	}
}
