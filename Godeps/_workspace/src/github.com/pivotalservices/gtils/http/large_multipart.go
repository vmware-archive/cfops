package http

import (
	"crypto/tls"
	"io"
	"net/http"

	"github.com/technoweenie/multipartstreamer"
	"github.com/xchapter7x/lo"
)

var LargeMultiPartUpload = func(conn ConnAuth, paramName, filename string, fileSize int64, fileRef io.Reader, params map[string]string) (res *http.Response, err error) {
	var req *http.Request
	ms := multipartstreamer.New()

	if params != nil {
		lo.G.Debug("adding params for request: ", params)

		if err = ms.WriteFields(params); err != nil {
			lo.G.Error("writefields error: ", err)
		}
	}

	if fileSize == int64(-1) {
		fileSize = GetFileSize(filename)
	}

	ms.WriteReader(paramName, filename, fileSize, fileRef)
	if req, err = http.NewRequest("POST", conn.Url, nil); err == nil {

		if conn.BearerToken != "" {
			req.Header.Add("Authorization", "Bearer "+conn.BearerToken)

		} else if conn.Username != "" && conn.Password != "" {
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
