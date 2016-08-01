package system

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"code.cloudfoundry.org/lager"
)

//Client ...
type Client interface {
	GetStagedProducts() ([]stagedProduct, error)
}

type opsmanClient struct {
	httpClient    httpClient
	token         string
	opsManagerURL string
	logger        lager.Logger
}

type stagedProduct struct {
	Type string `json:"type"`
	GUID string `json:"GUID"`
}

//NewOpsManagerClient ...
func NewOpsManagerClient(hostname, username, password string, logger lager.Logger) (Client, error) {
	url := "https://" + hostname

	token, err := getAuthToken(url, username, password)
	if err != nil {
		return &opsmanClient{}, err
	}

	opsManClient := newHttpClient(url, token, logger)
	return &opsmanClient{
		token:         token,
		httpClient:    opsManClient,
		opsManagerURL: url,
		logger:        logger,
	}, nil
}

func (client *opsmanClient) GetStagedProducts() ([]stagedProduct, error) {
	client.httpClient.NewRequest("api/v0/staged/products", nil)
	stagedProducts := []stagedProduct{}
	err := client.httpClient.Get(&stagedProducts)
	return stagedProducts, err
}

func getAuthToken(opsManagerURL, username, password string) (string, error) {
	body := url.Values{
		"grant_type": {"password"},
		"username":   {username},
		"password":   {password},
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	opsManClient := &http.Client{Transport: tr}
	request, err := http.NewRequest("POST", fmt.Sprintf("%s/uaa/oauth/token", opsManagerURL), strings.NewReader(body.Encode()))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	request.SetBasicAuth("opsman", "")
	response, err := opsManClient.Do(request)
	if err != nil {
		return "", err
	}
	responseToken := token{}

	if response.StatusCode != 200 {
		return "", fmt.Errorf("Unexpected response code %d", response.StatusCode)
	}

	err = json.NewDecoder(response.Body).Decode(&responseToken)
	if err != nil {
		return "", err
	}
	if responseToken.AccessToken == "" {
		return "", fmt.Errorf("No token returned")
	}

	return responseToken.AccessToken, nil
}

type token struct {
	AccessToken string `json:"access_token,required"`
}
