package apiService_repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jasanya-tech/jasanya-response-backend-golang/_error"
	"github.com/jasanya-tech/jasanya-response-backend-golang/response"
	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

func (m *ApiServiceRepositoryImpl) CreateProfileDueit(data []byte, appID string) (*repository.ProfileDueit, error) {
	endPoint := fmt.Sprintf("%s/account/profile/%s", m.baseUrlDueit, appID)

	dataReq := bytes.NewReader(data)
	req, err := http.NewRequest("POST", endPoint, dataReq)
	if err != nil {
		log.Warn().Msgf(util.LogErrHttpNewRequest, err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", infra.AppApiKeyAccount)

	client := http.Client{
		Timeout: 2 * time.Second,
	}

	responseReq, err := client.Do(req)
	if err != nil {
		log.Warn().Msgf(util.LogErrClientDo, err)
		err = _error.HttpErrString(response.CodeCompanyName[response.CM09], response.CM09)
		return nil, err
	}
	defer func() {
		if errBody := responseReq.Body.Close(); errBody != nil {
			log.Warn().Msgf(util.LogErrClientDoClose, err)
		}
	}()

	profile, err := m.fetchResponseProfile(responseReq)

	return profile, err
}

func (m *ApiServiceRepositoryImpl) GetProfileByUserIDDueit(userID string, appID string) (*repository.ProfileDueit, error) {
	endPoint := fmt.Sprintf("%s/account/profile/%s", m.baseUrlDueit, appID)

	req, err := http.NewRequest("GET", endPoint, nil)
	if err != nil {
		log.Warn().Msgf(util.LogErrHttpNewRequest, err)
		return nil, err
	}

	req.Header.Set("User-Id", userID)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", infra.AppApiKeyAccount)

	client := http.Client{
		Timeout: 2 * time.Second,
	}

	responseReq, err := client.Do(req)
	if err != nil {
		log.Warn().Msgf(util.LogErrClientDo, err)
		err = _error.HttpErrString(response.CodeCompanyName[response.CM09], response.CM09)
		return nil, err
	}
	defer func() {
		if errBody := responseReq.Body.Close(); errBody != nil {
			log.Warn().Msgf(util.LogErrClientDoClose, err)
		}
	}()

	profile, err := m.fetchResponseProfile(responseReq)

	return profile, err
}

func (m *ApiServiceRepositoryImpl) fetchResponseProfile(r *http.Response) (*repository.ProfileDueit, error) {
	if r.StatusCode == 200 {
		var profileRespMap map[string]json.RawMessage
		err := json.NewDecoder(r.Body).Decode(&profileRespMap)
		if err != nil {
			log.Warn().Msgf(util.LogErrDecode, r.Body, err)
			return nil, err
		}

		var profile repository.ProfileDueit
		err = json.Unmarshal(profileRespMap["data"], &profile)
		if err != nil {
			log.Warn().Msgf(util.LogErrUnmarshal, profileRespMap["data"], err)
			return nil, err
		}

		return &profile, nil
	}

	resp := map[string]any{}
	err := json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		log.Warn().Msgf(util.LogErrDecode, r.Body, err)
		return nil, err
	}

	log.Warn().Msgf("error fetch data in account service | response : %v", resp)

	statusString, ok := resp["status"].(string)
	if !ok {
		err = &response.HttpError{
			Err:         "bad gateway",
			CodeCompany: response.CM09,
		}
	}

	err = &response.HttpError{
		Err:         resp["errors"],
		CodeCompany: response.CodeCompany(statusString),
	}
	return nil, err
}
