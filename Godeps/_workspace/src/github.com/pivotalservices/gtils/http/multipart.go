package http

import (
	"bytes"
	"crypto/tls"
	"io"
	"mime/multipart"
	"net/http"
)

type MultiPartBodyFunc func(string, string, io.Reader, map[string]string) (io.Reader, error)

type UploadFunc func(string, string, string, string, string, io.Reader, map[string]string) (*http.Response, error)

func MultiPartBody(paramName, filename string, fileRef io.Reader, params map[string]string) (io.Reader, error) {
	var part io.Writer

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(paramName, filename)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, fileRef); err == nil {

		for key, val := range params {
			_ = writer.WriteField(key, val)
		}
		err = writer.Close()
		if err != nil {
			return nil, err
		}
	}
	return body, nil
}

func Upload(url, username, password, paramName, filename string, fileRef io.Reader, params map[string]string) (*http.Response, error) {
	var part io.Writer

	var body bytes.Buffer
	w := multipart.NewWriter(&body)

	part, err := w.CreateFormFile(paramName, filename)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, fileRef); err == nil {

		for key, val := range params {
			_ = w.WriteField(key, val)
		}
		err = w.Close()
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest("POST", url, &body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", w.FormDataContentType())
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	return client.Do(req)
}
