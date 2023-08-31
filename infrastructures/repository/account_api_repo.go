package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infrastructures/config"
	"github.com/rs/zerolog/log"
	"net/http"
)

type AccountApiRepoImpl struct {
}

func NewAccountApiRepoImpl() repository.AccountApiRepo {
	return &AccountApiRepoImpl{}
}

func (a *AccountApiRepoImpl) CreateProfile(ctx context.Context, data []byte) (*model.Profile, error) {
	endPoint := fmt.Sprintf("%s/profile", config.AppAccountApi)

	dataReq := bytes.NewReader(data)
	req, err := http.NewRequestWithContext(ctx, "POST", endPoint, dataReq)
	if err != nil {
		log.Err(err).Msg("failed request post to account service")
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Err(err).Msg("failed get response from http request post")
		return nil, err
	}
	defer func() {
		if errBody := response.Body.Close(); errBody != nil {
			log.Err(errBody).Msg("failed close response body")
		}
	}()

	var profile model.Profile

	err = json.NewDecoder(response.Body).Decode(&profile)
	log.Debug().Msg("process decoder response body")
	if err != nil {
		log.Err(err).Msg("failed decode response to struct")
		return nil, err
	}

	profile.Code = response.StatusCode

	return &profile, nil
}
