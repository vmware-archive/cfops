package http

import (
	"bytes"
	"crypto/tls"
	"io"
	"mime/multipart"
	"net/http"
)

type ConnAuth struct {
	Url         string
	Username    string
	Password    string
	BearerToken string
}

type MultiPartBodyFunc func(string, string, io.Reader, map[string]string) (io.ReadWriter, string, error)
type MultiPartUploadFunc func(conn ConnAuth, paramName, filename string, fileSize int64, fileRef io.Reader, params map[string]string) (res *http.Response, err error)

func MultiPartBody(paramName, filename string, fileRef io.Reader, params map[string]string) (body io.ReadWriter, contentType string, err error) {
	var part io.Writer

	body = &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if part, err = writer.CreateFormFile(paramName, filename); err == nil {

		if _, err = io.Copy(part, fileRef); err == nil {

			for key, val := range params {
				_ = writer.WriteField(key, val)
			}
			contentType = writer.FormDataContentType()
			writer.Close()
		}
	}
	return
}

var MultiPartUpload = func(conn ConnAuth, paramName, filename string, fileSize int64, fileRef io.Reader, params map[string]string) (res *http.Response, err error) {
	var contentType string
	var rbody io.Reader

	if rbody, contentType, err = MultiPartBody(paramName, filename, fileRef, params); err == nil {
		var req *http.Request

		if req, err = http.NewRequest("POST", conn.Url, rbody); err == nil {

			if conn.BearerToken != "" {
				req.Header.Add("Authorization", "Bearer "+conn.BearerToken)

			} else if conn.Username != "" && conn.Password != "" {
				req.SetBasicAuth(conn.Username, conn.Password)
			}
			req.Header.Set("Content-Type", contentType)
			client := NewTransportClient()
			res, err = client.Do(req)
		}
	}
	return
}

var NewTransportClient = func() (client interface {
	Do(*http.Request) (*http.Response, error)
}) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{Transport: tr}
	return
}
