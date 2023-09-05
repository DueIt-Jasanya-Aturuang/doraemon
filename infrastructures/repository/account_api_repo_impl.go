package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"
)

type AccountApiRepoImpl struct {
	endPoint string
}

func NewAccountApiRepoImpl(endPoint string) repository.AccountApiRepo {
	return &AccountApiRepoImpl{
		endPoint: endPoint,
	}
}

func (a *AccountApiRepoImpl) CreateProfile(data []byte) (*model.Profile, error) {
	endPoint := fmt.Sprintf("%s/profile", a.endPoint)

	dataReq := bytes.NewReader(data)
	req, err := http.NewRequest("POST", endPoint, dataReq)
	if err != nil {
		log.Err(err).Msg("failed request post to account service")
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 2 * time.Second,
	}
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

	if response.StatusCode == 201 {
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

	return nil, errors.New("BAD GATEWAY")
}

func (a *AccountApiRepoImpl) GetProfileByUserID(userID string) (*model.Profile, error) {
	endPoint := fmt.Sprintf("%s/profile", a.endPoint)

	req, err := http.NewRequest("GET", endPoint, nil)
	if err != nil {
		log.Err(err).Msg("failed request post to account service")
	}
	req.Header.Set("User-Id", userID)

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 2 * time.Second,
	}
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

	if response.StatusCode == 200 {
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

	return nil, errors.New("BAD GATEWAY")
}
