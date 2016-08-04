package system

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"code.cloudfoundry.org/lager"
)

type httpClient struct {
	baseURL       string
	authorization string
	client        *http.Client
	request       *http.Request
	logger        lager.Logger
}

//newHttpClient ...
func newHttpClient(url, authorization string, logger lager.Logger) httpClient {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return httpClient{
		baseURL:       url,
		authorization: authorization,
		client:        &http.Client{Transport: tr},
		logger:        logger.Session("http"),
	}
}

func (client *httpClient) NewRequest(endpoint string, body io.Reader) (*http.Request, error) {
	parsedURL, err := url.Parse(client.baseURL)
	if err != nil {
		return nil, err
	}
	parsedURL.Path = path.Join(parsedURL.Path, endpoint)
	request, err := http.NewRequest("", parsedURL.String(), body)
	if err != nil {
		panic(err)
	}
	client.request = request
	client.addCommonHeaders()
	return request, nil
}

func (client *httpClient) Post() error {
	client.request.Method = "POST"
	_, err := client.processRequest()
	return err
}

func (client *httpClient) Delete() error {
	client.request.Method = "DELETE"
	_, err := client.processRequest()
	return err
}

func (client *httpClient) Get(responseObject interface{}) error {
	client.request.Method = "GET"
	responseBody, err := client.processRequest()
	if err != nil {
		return err
	}
	err = json.NewDecoder(responseBody).Decode(&responseObject)
	// TODO: Add logging again, it fails case we have a json.RawMessage in a struct
	// TODO: uncomment when https://github.com/cloudfoundry/lager/pull/20 is fixed
	// client.logger.Debug("Response body", lager.Data{"response": responseObject})
	return err
}

func (client *httpClient) processRequest() (io.ReadCloser, error) {
	client.logger.Debug("making request", lager.Data{
		"URL":    client.request.URL.String(),
		"Method": client.request.Method,
	})
	response, err := client.client.Do(client.request)
	if err != nil {
		client.logger.Error("request failed", err)
		panic("Catz!")
	}

	if response.StatusCode != 200 {
		responseBody := map[string]interface{}{}
		json.NewDecoder(response.Body).Decode(&responseBody)

		return nil, fmt.Errorf("Request failed with %d. Info: %v", response.StatusCode, responseBody)
	}
	return response.Body, err
}

func (client *httpClient) addCommonHeaders() {
	client.request.Header.Set("Content-Type", "application/json")
	client.request.Header.Set("Authorization", client.authorization)
}
