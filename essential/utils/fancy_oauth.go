package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type FancyOauthConfig struct {
	ClientId        string
	ClientSecretKey string
	State           string
	RedirectUri     string
}

type AccessTokenResponse struct {
	RetCode int    `json:"retCode"`
	ReMsg   string `json:"reMsg"`
	RetData struct {
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
		AccessToken      string `json:"access_token"`
	} `json:"retData"`
}

type OauthUser struct {
	UserName string
}

type AuthInfoResponse struct {
	RetCode int    `json:"retCode"`
	ReMsg   string `json:"reMsg"`
	RetData struct {
		User *OauthUser
	} `json:"retData"`
}

const (
	LOGIN_URL            = "http://oa.fancydigital.com.cn/api/oauth/login/"
	LOGOUT_URL           = "http://oa.fancydigital.com.cn/api/oauth/logout/"
	GET_ACCESS_TOKEN_URL = "http://oa.fancydigital.com.cn/api/oauth/access-token"
	GET_USER_INFO_URL    = "http://oa.fancydigital.com.cn/api/oauth/get-user"
)

func NewConfig(clientId, clientSecretKey, state, redirectUri string) *FancyOauthConfig {
	return &FancyOauthConfig{
		ClientId:        clientId,
		ClientSecretKey: clientSecretKey,
		State:           state,
		RedirectUri:     redirectUri,
	}
}

func (c *FancyOauthConfig) FancyOauthLogin() string {
	params := &url.Values{}

	params.Set("client_id", c.ClientId)
	params.Set("redirect_uri", c.RedirectUri)
	params.Set("state", c.State)

	return LOGIN_URL + "?" + params.Encode()
}

func (c *FancyOauthConfig) GetAccessToken(code string) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", GET_ACCESS_TOKEN_URL, strings.NewReader("code="+code))

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+base64.URLEncoding.EncodeToString([]byte(c.ClientId+":"+c.ClientSecretKey)))

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	tokenData := AccessTokenResponse{}
	if err := json.Unmarshal(body, &tokenData); err != nil {
		return "", err
	}

	if tokenData.RetCode == 0 {
		return tokenData.RetData.AccessToken, nil
	}

	return "", fmt.Errorf("access token get fail: %s", tokenData.RetData.ErrorDescription)
}

func (c *FancyOauthConfig) GetAuthInfo(accessToken string) (*OauthUser, error) {
	resp, err := http.Post(GET_USER_INFO_URL,
		"application/x-www-form-urlencoded",
		strings.NewReader("access_token="+accessToken))

	if err != nil {
		return nil, fmt.Errorf("auth info request fail")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}
	authData := AuthInfoResponse{}
	if err := json.Unmarshal(body, &authData); err != nil {
		return nil, fmt.Errorf("auth info decode fail")
	}

	if authData.RetData.User != nil {
		return authData.RetData.User, nil
	}

	return nil, fmt.Errorf("auth info not find")
}
