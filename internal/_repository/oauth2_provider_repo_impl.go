package _repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/util/msg"
)

type GoogleOauthRepoImpl struct {
	clientID    string
	secretID    string
	redirectURI string
}

func NewGoogleOauthRepoImpl(
	clientID string,
	secretID string,
	redirectURI string,
) repository.Oauth2ProviderRepo {
	return &GoogleOauthRepoImpl{
		clientID:    clientID,
		secretID:    secretID,
		redirectURI: redirectURI,
	}
}

func (g *GoogleOauthRepoImpl) GetGoogleOauthToken(code string) (*model.GoogleOauth2Token, error) {
	const uri = "https://oauth2.googleapis.com/token"

	value := url.Values{}
	value.Add("grant_type", "authorization_code")
	value.Add("code", code)
	value.Add("client_id", g.clientID)
	value.Add("secret_id", g.secretID)
	value.Add("redirect_uri", g.redirectURI)

	query := value.Encode()
	queryBuffer := bytes.NewBufferString(query)

	req, err := http.NewRequest("POST", uri, queryBuffer)
	if err != nil {
		log.Err(err).Msg(_msg.LogErrHttpNewRequest)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{
		Timeout: 2 * time.Second,
	}

	response, err := client.Do(req)
	if err != nil {
		log.Err(err).Msg(_msg.LogErrHttpClientDo)
		return nil, err
	}
	defer func() {
		if errBody := response.Body.Close(); errBody != nil {
			log.Err(errBody).Msg(_msg.LogErrResponseBodyClose)
		}
	}()

	log.Debug().Msg(response.Status)
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("could not retrieve token")
	}

	var respBody bytes.Buffer
	_, err = io.Copy(&respBody, response.Body)
	if err != nil {
		log.Err(err).Msg("failed copy response body to bytes buffer")
		return nil, err
	}

	var googleOauth2TokenMap map[string]any

	err = json.Unmarshal(respBody.Bytes(), &googleOauth2TokenMap)
	if err != nil {
		log.Err(err).Msg("failed unmarshal response body bytes buffer to map")
		return nil, err
	}

	googleOauthToken := model.GoogleOauth2Token{
		AccessToken: googleOauth2TokenMap["access_token"].(string),
		IDToken:     googleOauth2TokenMap["id_token"].(string),
	}

	return &googleOauthToken, nil
}

func (g *GoogleOauthRepoImpl) GetGoogleOauthUser(token *model.GoogleOauth2Token) (*model.GoogleOauth2User, error) {
	uri := fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?alt=json&access_token=%s", token.AccessToken)

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Err(err).Msg(_msg.LogErrHttpNewRequest)
		return nil, err
	}

	req.Header.Set("Authorization", token.IDToken)

	client := http.Client{
		Timeout: 2 * time.Second,
	}

	response, err := client.Do(req)
	if err != nil {
		log.Err(err).Msg(_msg.LogErrHttpClientDo)
		return nil, err
	}
	defer func() {
		if errBody := response.Body.Close(); errBody != nil {
			log.Err(errBody).Msg(_msg.LogErrResponseBodyClose)
		}
	}()

	if response.StatusCode != http.StatusOK {
		log.Warn().Msg("status code not 200")
		return nil, errors.New("could not retrieve user")
	}

	var respBody bytes.Buffer
	_, err = io.Copy(&respBody, response.Body)
	if err != nil {
		log.Err(err).Msg("failed copy response body to bytes buffer")
		return nil, err
	}

	var googleOauth2UserMap map[string]any

	err = json.Unmarshal(respBody.Bytes(), &googleOauth2UserMap)
	if err != nil {
		log.Err(err).Msg("failed unmarshal response body bytes buffer to map")
		return nil, err
	}

	userBody := &model.GoogleOauth2User{
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
