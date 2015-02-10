package http

import (
	"bytes"
	"io"
	"mime/multipart"
)

type MultiPartBodyFunc func(string, string, io.Reader, map[string]string) (io.Reader, error)

func MultiPartBody(paramName, filename string, fileRef io.Reader, params map[string]string) (body io.Reader, err error) {
	var part io.Writer

	bodyBuffer := &bytes.Buffer{}
	body = bodyBuffer
	writer := multipart.NewWriter(bodyBuffer)

	if part, err = writer.CreateFormFile(paramName, filename); err == nil {

		if _, err = io.Copy(part, fileRef); err == nil {

			for key, val := range params {
				_ = writer.WriteField(key, val)
			}
			writer.Close()
		}
	}
	return
}
