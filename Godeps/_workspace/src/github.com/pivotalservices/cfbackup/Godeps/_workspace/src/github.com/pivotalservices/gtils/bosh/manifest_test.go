package bosh

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"gopkg.in/yaml.v2"
)

func ReadFakeResponse(statusCode int) (res *http.Response) {
	file, _ := os.Open("fixtures/manifest.yml")
	response := &http.Response{}
	response.Body = file
	response.StatusCode = statusCode
	return response
}

type FakeManifest struct {
	Compilation string
	Job         string
}

func TestRetriveManifestFromHttpResponse(t *testing.T) {
	resp := ReadFakeResponse(200)
	reader, err := retrieveManifest(resp)
	if err != nil {
		t.Fatalf("Failed to retrieve Manifest : %v", err)
	}
	m := FakeManifest{}
	bytes, err := ioutil.ReadAll(reader.(io.Reader))
	err = yaml.Unmarshal(bytes, &m)
	if err != nil {
		t.Fatalf("Failed to parse the  Manifest : %v", err)
	}
	if m.Job != "test" || m.Compilation != "test" {
		t.Errorf("The retrieved manifest not valid: %v", m)
	}
}
