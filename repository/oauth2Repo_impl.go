package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/jasanya-tech/jasanya-response-backend-golang/_error"
	"github.com/jasanya-tech/jasanya-response-backend-golang/response"
	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

type Oauth2RepositoryImpl struct {
	clientID    string
	secretID    string
	redirectURI string
}

func NewOauth2RepositoryImpl(
	clientID string,
	secretID string,
	redirectURI string,
) domain.Oauth2Repository {
	return &Oauth2RepositoryImpl{
		clientID:    clientID,
		secretID:    secretID,
		redirectURI: redirectURI,
	}
}

func (o *Oauth2RepositoryImpl) GetGoogleToken(code string) (*domain.Oauth2GoogleToken, error) {
	const uri = "https://oauth2.googleapis.com/token"

	value := url.Values{}
	value.Add("grant_type", "authorization_code")
	value.Add("code", code)
	value.Add("client_id", o.clientID)
	value.Add("secret_id", o.secretID)
	value.Add("redirect_uri", o.redirectURI)

	query := value.Encode()
	queryBuffer := bytes.NewBufferString(query)

	req, err := http.NewRequest("POST", uri, queryBuffer)
	if err != nil {
		log.Warn().Msgf(util.LogErrHttpNewRequest, err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{
		Timeout: 2 * time.Second,
	}

	responseReq, err := client.Do(req)
	if err != nil {
		log.Warn().Msgf(util.LogErrClientDo, err)
		return nil, err
	}
	defer func() {
		if errBody := responseReq.Body.Close(); errBody != nil {
			log.Warn().Msgf(util.LogErrClientDoClose, err)
		}
	}()

	if responseReq.StatusCode != http.StatusOK {
		return nil, _error.HttpErrString("invalid token", response.CM05)
	}

	var respBody bytes.Buffer
	_, err = io.Copy(&respBody, responseReq.Body)
	if err != nil {
		log.Err(err).Msgf("failed copy responseReq body to bytes buffer | dst : %v | src : %v", respBody, responseReq.Body)
		return nil, err
	}

	var googleOauth2TokenMap map[string]any

	err = json.Unmarshal(respBody.Bytes(), &googleOauth2TokenMap)
	if err != nil {
		log.Warn().Msgf(util.LogErrUnmarshal, respBody.Bytes(), err)
		return nil, err
	}

	googleOauthToken := domain.Oauth2GoogleToken{
		AccessToken: googleOauth2TokenMap["access_token"].(string),
		IDToken:     googleOauth2TokenMap["id_token"].(string),
	}

	return &googleOauthToken, nil
}

func (o *Oauth2RepositoryImpl) GetGoogleUser(token *domain.Oauth2GoogleToken) (*domain.Oauth2GoogleUser, error) {
	uri := fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?alt=json&access_token=%s", token.AccessToken)

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Warn().Msgf(util.LogErrHttpNewRequest, err)
		return nil, err
	}

	req.Header.Set("Authorization", token.IDToken)

	client := http.Client{
		Timeout: 2 * time.Second,
	}

	responseReq, err := client.Do(req)
	if err != nil {
		log.Warn().Msgf(util.LogErrClientDo, err)
		return nil, err
	}
	defer func() {
		if errBody := responseReq.Body.Close(); errBody != nil {
			log.Warn().Msgf(util.LogErrClientDoClose, err)
		}
	}()

	if responseReq.StatusCode != http.StatusOK {
		log.Warn().Msgf("failed recive user | responseReq : %v", responseReq)
		return nil, _error.HttpErrString("invalid token", response.CM05)
	}

	var respBody bytes.Buffer
	_, err = io.Copy(&respBody, responseReq.Body)
	if err != nil {
		log.Err(err).Msgf("failed copy responseReq body to bytes buffer | dst : %v | src : %v", respBody, responseReq.Body)
		return nil, err
	}

	var googleOauth2UserMap map[string]any

	err = json.Unmarshal(respBody.Bytes(), &googleOauth2UserMap)
	if err != nil {
		log.Warn().Msgf(util.LogErrUnmarshal, respBody.Bytes(), err)
		return nil, err
	}

	userBody := &domain.Oauth2GoogleUser{
		ID:            googleOauth2UserMap["id"].(string),
		Email:         googleOauth2UserMap["email"].(string),
		VerifiedEmail: googleOauth2UserMap["verified_email"].(bool),
		Name:          googleOauth2UserMap["name"].(string),
		GivenName:     googleOauth2UserMap["given_name"].(string),
		Image:         googleOauth2UserMap["picture"].(string),
		Locale:        googleOauth2UserMap["locale"].(string),
	}

	return userBody, nil
}
