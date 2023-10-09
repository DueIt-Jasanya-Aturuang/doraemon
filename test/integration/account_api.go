package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/pkg/repository"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra"
)

func CreateProfile(t *testing.T) {
	infra.LogInit()
	infra.BaseUrlDueitAccountService = "http://localhost:8181"
	accountApi := repository.NewAccountApiRepoImpl(infra.BaseUrlDueitAccountService)

	req := []byte(`{
			"user_id": "123"
		 }`)

	profile, err := accountApi.CreateProfile(req)
	assert.Error(t, err)
	assert.Nil(t, profile)
}
