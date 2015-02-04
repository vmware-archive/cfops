package httptest

import (
	"io"

	"github.com/pivotalservices/gtils/http"
)

type EntityCaptcherFunc func(http.HttpRequestEntity)

type MockGateway struct {
	FakeGetAdaptor  http.RequestAdaptor
	FakePutAdaptor  http.RequestAdaptor
	FakePostAdaptor http.RequestAdaptor
	Capture         EntityCaptcherFunc
}

func (gateway *MockGateway) Get(entity http.HttpRequestEntity) http.RequestAdaptor {
	gateway.Capture(entity)
	return gateway.FakeGetAdaptor
}

func (gateway *MockGateway) Put(entity http.HttpRequestEntity, body io.Reader) http.RequestAdaptor {
	gateway.Capture(entity)
	return gateway.FakePutAdaptor
}

func (gateway *MockGateway) Post(entity http.HttpRequestEntity, body io.Reader) http.RequestAdaptor {
	gateway.Capture(entity)
	return gateway.FakePostAdaptor
}
