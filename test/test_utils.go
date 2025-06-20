package test

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type Builder struct {
	requestPath string
	statusCode  int
	body        interface{}
}

func (b *Builder) Execute() {
	interceptor := &Interceptor{}
	interceptor.Add(b.requestPath, b.statusCode, b.body)
	http.DefaultClient = &http.Client{Transport: interceptor}
}

func (b *Builder) ReturnStatusCode(statusCode int) *Builder {
	b.statusCode = statusCode
	return b
}

func (b *Builder) ReturnBody(body interface{}) *Builder {
	b.body = body
	return b
}

func When(requestPath string) *Builder {
	b := &Builder{}
	b.requestPath = requestPath
	return b
}

type Interceptor struct {
	Transport http.RoundTripper
	option    Option
}

type Option struct {
	requestPath string
	statusCode  int
	body        interface{}
}

func (interceptor *Interceptor) Add(requestPath string, statusCode int, body interface{}) {
	interceptor.option = Option{
		requestPath: requestPath,
		statusCode:  statusCode,
		body:        body,
	}
	interceptor.Transport = http.DefaultTransport
	http.DefaultClient = &http.Client{Transport: interceptor}
}

func (interceptor *Interceptor) RoundTrip(req *http.Request) (*http.Response, error) {

	validResponse, _ := json.Marshal(interceptor.option.body)

	return &http.Response{
		StatusCode: interceptor.option.statusCode,
		Body:       io.NopCloser(strings.NewReader(string(validResponse))),
	}, nil

}
