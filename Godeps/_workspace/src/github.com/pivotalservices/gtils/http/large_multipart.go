package http

import (
	"io"
	"net/http"

	"github.com/technoweenie/multipartstreamer"
)

var LargeMultiPartUpload = func(conn ConnAuth, paramName, filename string, fileRef io.Reader, params map[string]string) (res *http.Response, err error) {
	ms := multipartstreamer.New()
	ms.WriteFile(paramName, filename)
	if req, err := http.NewRequest("POST", conn.Url, nil); err == nil {
		if conn.Username != "" && conn.Password != "" {
			req.SetBasicAuth(conn.Username, conn.Password)
		}
		ms.SetupRequest(req)
		res, err = http.DefaultClient.Do(req)
	}
	return
}
