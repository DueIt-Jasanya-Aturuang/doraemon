package integration

import (
	"context"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infrastructures/config"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infrastructures/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateProfile(t *testing.T) {
	config.LogInit()
	ctx := context.Background()
	config.AppAccountApi = "http://localhost:8181"
	accountApi := repository.NewAccountApiRepoImpl()

	req := []byte(`{
			"user_id": "123"
		 }`)

	profile, err := accountApi.CreateProfile(ctx, req)
	t.Log(profile)
	assert.NoError(t, err)
	assert.NotNil(t, profile)
}
