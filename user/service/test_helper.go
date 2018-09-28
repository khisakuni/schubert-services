package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testCase struct {
	before func(*testing.T)
	params interface{}
	code   int
	body   string
}

func newRequest(method, endpoint string, p interface{}) (*http.Request, error) {
	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(p)
	if err != nil {
		return nil, err
	}
	return http.NewRequest(method, endpoint, body)
}

func sendRequest(t *testing.T, service Service, req *http.Request) *httptest.ResponseRecorder {
	res := httptest.NewRecorder()
	service.Router.ServeHTTP(res, req)
	return res
}
