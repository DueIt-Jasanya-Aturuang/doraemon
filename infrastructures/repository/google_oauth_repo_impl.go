package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"net/url"
	"time"
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
) repository.GoogleOauthRepo {
	return &GoogleOauthRepoImpl{
		clientID:    clientID,
		secretID:    secretID,
		redirectURI: redirectURI,
	}
}

func (g *GoogleOauthRepoImpl) GetGoogleOauthToken(code string) (*model.GoogleOauthToken, error) {
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
		log.Err(err).Msg("failed request to oauth token")
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{
		Timeout: 2 * time.Second,
	}

	response, err := client.Do(req)
	if err != nil {
		log.Err(err).Msg("failed to load response from request post")
		return nil, err
	}
	defer func() {
		if errBody := response.Body.Close(); errBody != nil {
			log.Err(errBody).Msg("failed close response body")
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

	var googleOauthTokenMap map[string]any

	err = json.Unmarshal(respBody.Bytes(), &googleOauthTokenMap)
	if err != nil {
		log.Err(err).Msg("failed unmarshal response body bytes buffer to map")
		return nil, err
	}

	googleOauthToken := model.GoogleOauthToken{
		AccessToken: googleOauthTokenMap["access_token"].(string),
		IDToken:     googleOauthTokenMap["id_token"].(string),
	}

	return &googleOauthToken, nil
}

func (g *GoogleOauthRepoImpl) GetGoogleOauthUser(token *model.GoogleOauthToken) (*model.GoogleOauthUser, error) {
	uri := fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?alt=json&access_token=%s", token.AccessToken)

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Err(err).Msg("failed request to google apis method get")
		return nil, err
	}

	req.Header.Set("Authorization", token.IDToken)

	client := http.Client{
		Timeout: 2 * time.Second,
	}

	response, err := client.Do(req)
	if err != nil {
		log.Err(err).Msg("failed receive response from request")
		return nil, err
	}
	defer func() {
		if errBody := response.Body.Close(); errBody != nil {
			log.Err(errBody).Msg("failed close response body")
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

	var googleOauthUserMap map[string]any

	err = json.Unmarshal(respBody.Bytes(), &googleOauthUserMap)
	if err != nil {
		log.Err(err).Msg("failed unmarshal response body bytes buffer to map")
		return nil, err
	}

	userBody := &model.GoogleOauthUser{
		ID:            googleOauthUserMap["id"].(string),
		Email:         googleOauthUserMap["email"].(string),
		VerifiedEmail: googleOauthUserMap["verified_email"].(bool),
		Name:          googleOauthUserMap["name"].(string),
		GivenName:     googleOauthUserMap["given_name"].(string),
		Image:         googleOauthUserMap["picture"].(string),
		Locale:        googleOauthUserMap["locale"].(string),
	}

	return userBody, nil
}
