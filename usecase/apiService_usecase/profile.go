package apiService_usecase

import (
	"encoding/json"
	"errors"

	"github.com/jasanya-tech/jasanya-response-backend-golang/response"
	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

func (a *ApiServiceUsecaseImpl) GetProfileByUserID(userID string, appID string) (*usecase.ResponseProfileDueit, error) {
	profile, err := a.apiServiceRepo.GetProfileByUserIDDueit(userID, appID)
	if err != nil {
		var errhttp *response.HttpError
		if errors.As(err, &errhttp) {
			if errhttp.CodeCompany == response.CM01 {
				profileReq := usecase.RequestCreateProfile{
					UserID: userID,
				}

				profileJson, err := json.Marshal(&profileReq)
				if err != nil {
					log.Warn().Msgf(util.LogErrMarshal, userID, err)
					return nil, err
				}

				profile, err = a.apiServiceRepo.CreateProfileDueit(profileJson, appID)
				if err != nil {
					return nil, err
				}

				return &usecase.ResponseProfileDueit{
					ProfileID: profile.ProfileID,
					Quote:     profile.Quote,
					Profesi:   profile.Profesi,
				}, nil
			}
		}

		return nil, err
	}

	return &usecase.ResponseProfileDueit{
		ProfileID: profile.ProfileID,
		Quote:     profile.Quote,
		Profesi:   profile.Profesi,
	}, nil
}
