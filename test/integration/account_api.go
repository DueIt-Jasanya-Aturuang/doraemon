package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra/config"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/_repository"
)

func CreateProfile(t *testing.T) {
	config.LogInit()
	config.AppAccountApi = "http://localhost:8181"
	accountApi := _repository.NewAccountApiRepoImpl(config.AppAccountApi)

	req := []byte(`{
			"user_id": "123"
		 }`)

	profile, err := accountApi.CreateProfile(req)
	assert.Error(t, err)
	assert.Nil(t, profile)
}
