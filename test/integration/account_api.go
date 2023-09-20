package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/_repository"
)

func CreateProfile(t *testing.T) {
	infra.LogInit()
	infra.AppAccountApi = "http://localhost:8181"
	accountApi := _repository.NewAccountApiRepoImpl(infra.AppAccountApi)

	req := []byte(`{
			"user_id": "123"
		 }`)

	profile, err := accountApi.CreateProfile(req)
	assert.Error(t, err)
	assert.Nil(t, profile)
}
