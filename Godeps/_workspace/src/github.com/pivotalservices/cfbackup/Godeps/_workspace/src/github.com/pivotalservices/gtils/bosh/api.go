package bosh

import (
	"fmt"
	"net/url"
	"regexp"

	. "github.com/pivotalservices/gtils/http"
)

type API struct {
	Path           string
	ContentType    string
	Method         string
	HandleResponse HandleRespFunc
}

type regexReplace func(string) string

func replaceUrlParams(params map[string]string) regexReplace {
	return func(str string) string {
		var retValue string
		for key, value := range params {
			matchString := fmt.Sprintf("{%s}", key)
			if str == matchString {
				retValue = value
			}
		}
		return retValue
	}
}

var ParseUrl = func(ip string, port int, pathPattern string, pathParams, queryParams map[string]string) (urlString string, err error) {
	host := fmt.Sprintf("https://%s:%d/%s", ip, port, pathPattern)
	urlString = regexp.MustCompile("{(.+?)}").ReplaceAllStringFunc(host, replaceUrlParams(pathParams))
	u, err := url.Parse(urlString)
	if err != nil {
		return
	}
	q := u.Query()
	for key, value := range queryParams {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}
