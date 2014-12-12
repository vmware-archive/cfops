package http

import (
	"net/http"
)

type Handler interface {
	Handle(response *http.Response) (interface{}, error)
}
