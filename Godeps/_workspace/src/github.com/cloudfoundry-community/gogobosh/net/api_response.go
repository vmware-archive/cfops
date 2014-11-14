package net

import (
	"fmt"
	"net/http"
)

type ApiResponse struct {
	Message    string
	ErrorCode  string
	StatusCode int
	RedirectLocation string

	isError    bool
	isNotFound bool
	isRedirection bool
}

func NewApiResponse(message string, errorCode string, statusCode int) (apiResponse ApiResponse) {
	return ApiResponse{
		Message:    message,
		ErrorCode:  errorCode,
		StatusCode: statusCode,
		isError:    true,
	}
}

func NewApiResponseWithRedirect(location string) (apiResponse ApiResponse) {
	return ApiResponse{
		StatusCode:       http.StatusFound,
		RedirectLocation: location,
		isRedirection:    true,
	}
}

func NewApiResponseWithStatusCode(statusCode int) (apiResponse ApiResponse) {
	return ApiResponse{
		StatusCode: statusCode,
	}
}

func NewApiResponseWithMessage(message string, a ...interface{}) (apiResponse ApiResponse) {
	return ApiResponse{
		Message: fmt.Sprintf(message, a...),
		isError: true,
	}
}

func NewApiResponseWithError(message string, err error) (apiResponse ApiResponse) {
	return ApiResponse{
		Message: fmt.Sprintf("%s: %s", message, err.Error()),
		isError: true,
	}
}

func NewNotFoundApiResponse(message string, a ...interface{}) (apiResponse ApiResponse) {
	return ApiResponse{
		Message:    fmt.Sprintf(message, a...),
		isNotFound: true,
	}
}

func NewSuccessfulApiResponse() (apiResponse ApiResponse) {
	return ApiResponse{}
}

func (apiResponse ApiResponse) IsError() bool {
	return apiResponse.isError
}

func (apiResponse ApiResponse) IsNotFound() bool {
	return apiResponse.isNotFound
}

func (apiResponse ApiResponse) IsRedirection() bool {
	return apiResponse.isRedirection
}

func (apiResponse ApiResponse) IsSuccessful() bool {
	return !apiResponse.IsNotSuccessful()
}

func (apiResponse ApiResponse) IsNotSuccessful() bool {
	return apiResponse.IsError() || apiResponse.IsNotFound()
}
