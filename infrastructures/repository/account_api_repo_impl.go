package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"
	_msg "github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/msg"
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
	endPoint := fmt.Sprintf("%s/account/profile", a.endPoint)

	dataReq := bytes.NewReader(data)
	req, err := http.NewRequest("POST", endPoint, dataReq)
	if err != nil {
		log.Err(err).Msg(_msg.LogErrHttpNewRequest)
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 2 * time.Second,
	}
	response, err := client.Do(req)
	if err != nil {
		var errHTTP = &model.ErrResponseHTTP{
			Code:    502,
			Message: "bad gateway",
		}
		return nil, errHTTP
	}
	defer func() {
		if errBody := response.Body.Close(); errBody != nil {
			log.Err(errBody).Msg(_msg.LogErrResponseBodyClose)
		}
	}()

	if response.StatusCode == 201 {
		var profileRespMap map[string]json.RawMessage
		err = json.NewDecoder(response.Body).Decode(&profileRespMap)
		if err != nil {
			log.Err(err).Msg(_msg.LogErrJsonNewDecoderDecode)
			return nil, err
		}

		var profile model.Profile
		err = json.Unmarshal(profileRespMap["data"], &profile)
		if err != nil {
			log.Err(err).Msg(_msg.LogErrJsonNewDecoderDecode)
			return nil, err
		}

		return &profile, nil
	} else {
		resp := map[string]any{}
		err = json.NewDecoder(response.Body).Decode(&resp)
		var errHTTP = &model.ErrResponseHTTP{
			Code:    response.StatusCode,
			Message: resp["errors"],
		}
		return nil, errHTTP
	}

}

func (a *AccountApiRepoImpl) GetProfileByUserID(userID string) (*model.Profile, error) {
	endPoint := fmt.Sprintf("%s/account/profile", a.endPoint)

	req, err := http.NewRequest("GET", endPoint, nil)
	if err != nil {
		log.Err(err).Msg(_msg.LogErrHttpNewRequest)
	}
	req.Header.Set("User-Id", userID)

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 2 * time.Second,
	}
	response, err := client.Do(req)
	if err != nil {
		var errHTTP = &model.ErrResponseHTTP{
			Code:    502,
			Message: "bad gateway",
		}
		return nil, errHTTP
	}
	defer func() {
		if errBody := response.Body.Close(); errBody != nil {
			log.Err(errBody).Msg(_msg.LogErrResponseBodyClose)
		}
	}()

	if response.StatusCode == 200 {
		var profileRespMap map[string]json.RawMessage
		err = json.NewDecoder(response.Body).Decode(&profileRespMap)
		if err != nil {
			log.Err(err).Msg(_msg.LogErrJsonNewDecoderDecode)
			return nil, err
		}

		var profile model.Profile
		err = json.Unmarshal(profileRespMap["data"], &profile)
		if err != nil {
			log.Err(err).Msg(_msg.LogErrJsonNewDecoderDecode)
			return nil, err
		}

		return &profile, nil
	} else {
		resp := map[string]any{}
		err = json.NewDecoder(response.Body).Decode(&resp)
		var errHTTP = &model.ErrResponseHTTP{
			Code:    response.StatusCode,
			Message: resp["errors"],
		}
		return nil, errHTTP
	}
}
