package fake

import (
	"io"
	"io/ioutil"
	"net/http"

	ghttp "github.com/pivotalservices/gtils/http"
)

type HttpRequestor struct {
	ghttp.HttpGateway
	GetCount  int
	PostCount int
	PutCount  int
}

func (s *HttpRequestor) Get(entity ghttp.HttpRequestEntity) ghttp.RequestAdaptor {
	return func() (x *http.Response, y error) {
		return
	}
}

func (s *HttpRequestor) Post(entity ghttp.HttpRequestEntity, body io.Reader) ghttp.RequestAdaptor {
	return func() (x *http.Response, y error) {
		return
	}
}

func (s *HttpRequestor) Put(entity ghttp.HttpRequestEntity, body io.Reader) ghttp.RequestAdaptor {
	return func() (x *http.Response, y error) {
		return
	}
}

type MultiPart struct {
	UploadCallCount int
	SpyFileContents []byte
	StatusCode      int
}

func (s *MultiPart) Upload(conn ghttp.ConnAuth, paramName, filename string, fileSize int64, fileRef io.Reader, params map[string]string) (res *http.Response, err error) {
	s.SpyFileContents, _ = ioutil.ReadAll(fileRef)
	s.UploadCallCount++
	res = &http.Response{
		StatusCode: s.StatusCode,
	}
	return
}
