package http

import (
	"net/http"
)

type HttpResponseHandler interface {
	Handle(response *http.Response) (interface{}, error)
}
