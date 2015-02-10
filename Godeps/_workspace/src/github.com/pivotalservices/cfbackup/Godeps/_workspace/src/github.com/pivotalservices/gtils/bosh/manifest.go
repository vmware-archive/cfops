package bosh

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type ManifestResponse struct {
	Manifest string `json:"manifest"`
}

func retrieveManifest(response *http.Response) (manifest io.Reader, err error) {
	if response.StatusCode != 200 {
		err = ManifestStatusCodeError
		return
	}
	m := ManifestResponse{}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &m)
	if err != nil {
		return
	}
	return strings.NewReader(m.Manifest), nil
}
