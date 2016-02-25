package uaa

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func GetToken(uaaURL, opsManagerUsername, opsManagerPassword, clientID, clientSecret string) (token string, err error) {
	var res *http.Response
	params := url.Values{
		"grant_type":    {"password"},
		"response_type": {"token"},
		"username":      {opsManagerUsername},
		"password":      {opsManagerPassword},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
	}
	tokenURL := fmt.Sprintf("%s/oauth/token", uaaURL)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	if res, err = client.PostForm(tokenURL, params); err == nil && res.StatusCode == http.StatusOK {
		var body []byte

		if body, err = ioutil.ReadAll(res.Body); err == nil {
			t := new(Token)

			if err = json.Unmarshal(body, &t); err == nil {
				token = t.AccessToken
			}
		}
	} else if res != nil && res.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		err = errors.New(string(body))
	}
	return
}
