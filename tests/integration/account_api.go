package integration

import (
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infrastructures/config"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infrastructures/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

func CreateProfile(t *testing.T) {
	config.LogInit()
	config.AppAccountApi = "http://localhost:8181"
	accountApi := repository.NewAccountApiRepoImpl(config.AppAccountApi)

	req := []byte(`{
			"user_id": "123"
		 }`)

	profile, err := accountApi.CreateProfile(req)
	//assert.NoError(t, err)
	assert.Error(t, err)
	//assert.NotNil(t, profile)
	assert.Nil(t, profile)
}
