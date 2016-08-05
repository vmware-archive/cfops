package system

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	. "github.com/onsi/ginkgo"

	"code.cloudfoundry.org/lager"
)

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

type credentials struct {
	Credential struct {
		Type  string `json:"simple_credentials"`
		Value struct {
			Username string `json:"identity"`
			Password string `json:"password"`
		}
	}
}

//NewOpsManagerClient ...
func NewOpsManagerClient(hostname, username, password string, logger lager.Logger) (*opsmanClient, error) {
	url := "https://" + hostname

	token, err := getAuthorization(url, username, password)

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

type installationSettings struct {
	Products products `json:"products"`
}
type products []product

func (products products) CF() product {
	for _, p := range products {
		if strings.HasPrefix(p.InstallationName, "cf-") {
			return p
		}
	}
	Fail("Cant find cf product in installation settings")
	return product{}
}

type product struct {
	InstallationName string `json:"installation_name"`
	Jobs             jobs   `json:"jobs"`
}
type jobs []job

func (jobs jobs) UAA() job {
	for _, j := range jobs {
		if j.InstallationName == "uaa" {
			return j
		}
	}
	Fail("Cant find uaa job in product settings")
	return job{}
}

type job struct {
	InstallationName string `json:"installation_name"`
	Properties       properties
}

type properties []property

func (properties properties) AdminCredentials() property {
	for _, p := range properties {
		if p.Identifier == "admin_credentials" {
			return p
		}
	}
	Fail("Cant find admin_credentials in uaa job")
	return property{}
}

type property struct {
	Identifier string          `json:"identifier"`
	Value      json.RawMessage `json:"value"`
}

type credentialsValue struct {
	Username string `json:"identity"`
	Password string `json:"password"`
}

func (property property) Credentials() (username, password string, err error) {
	creds := credentialsValue{}
	fmt.Printf("property.Value %s\n", string(property.Value))
	err = json.Unmarshal(property.Value, &creds)
	return creds.Username, creds.Password, err
}

func (client *opsmanClient) GetInstallationSettings() (installationSettings, error) {
	client.httpClient.NewRequest("api/installation_settings", nil)
	installationSetting := installationSettings{}
	err := client.httpClient.Get(&installationSetting)
	return installationSetting, err
}

func (client *opsmanClient) GetAdminCredentials() (username, password string, err error) {
	installationSetting, err := client.GetInstallationSettings()
	if err != nil {
		return "", "", err
	}
	fmt.Printf("installationSetting.Products.CF().Jobs %d\n", len(installationSetting.Products.CF().Jobs))
	return installationSetting.Products.CF().Jobs.UAA().Properties.AdminCredentials().Credentials()
}

func getAuthorization(opsManagerURL, username, password string) (string, error) {
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

	if response.StatusCode == 200 {
		responseToken := token{}
		err = json.NewDecoder(response.Body).Decode(&responseToken)
		if err != nil {
			return "", err
		}
		if responseToken.AccessToken == "" {
			return "", fmt.Errorf("No token returned")
		}
		return "Bearer " + responseToken.AccessToken, nil
	} else if response.StatusCode == 404 {
		// Assume basic auth
		auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
		return "Basic " + auth, nil
	} else {
		return "", fmt.Errorf("Unexpected response code %d", response.StatusCode)
	}
}

type token struct {
	AccessToken string `json:"access_token,required"`
}
