package main

import (
	"net/http"
	"strings"
)

type HttpRequest struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    string            `json:"body,omitempty"`
}

func (r *HttpRequest) ToHttpRequest() (*http.Request, error) {
	req, err := http.NewRequest(r.Method, r.URL, strings.NewReader(r.Body))
	if err != nil {
		return nil, err
	}
	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}
	return req, nil
}
