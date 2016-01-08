package http

import (
	"crypto/tls"
	"io"
	"net/http"

	"github.com/technoweenie/multipartstreamer"
)

var LargeMultiPartUpload = func(conn ConnAuth, paramName, filename string, fileSize int64, fileRef io.Reader, params map[string]string) (res *http.Response, err error) {
	var req *http.Request
	ms := multipartstreamer.New()
	ms.WriteReader(paramName, filename, fileSize, fileRef)
	if req, err = http.NewRequest("POST", conn.Url, nil); err == nil {
		if conn.Username != "" && conn.Password != "" {
			req.SetBasicAuth(conn.Username, conn.Password)
		}
		ms.SetupRequest(req)
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		res, err = client.Do(req)
	}
	return
}
