package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/go-chi/chi/v5"
)

type fn func() *chi.Mux

// TestAPICall allows an easy way to test HTTP end points in unit testing
func TestAPICall(method string, endpoint string, data io.Reader, handler http.HandlerFunc) (code int, body *bytes.Buffer, err error) {
	if !strings.HasPrefix(endpoint, "/") {
		endpoint = "/" + endpoint
	}
	req, err := http.NewRequest(method, endpoint, data)
	if err != nil {
		return 500, nil, err
	}

	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	rr := httptest.NewRecorder()

	chi := Setup()
	chi.ServeHTTP(rr, req)

	return rr.Code, rr.Body, nil
}

// UnmarshalTestMap helps to unmarshal the request for the testing calls
func UnmarshalTestMap(body *bytes.Buffer) (APIReturn, map[string]interface{}, error) {
	ret := APIReturn{}
	retBuf := new(bytes.Buffer)
	retBuf.ReadFrom(body)
	err := json.Unmarshal(retBuf.Bytes(), &ret)
	if err != nil {
		return ret, map[string]interface{}{}, err
	}
	retBody, ok := ret.Data.(map[string]interface{})
	if !ok {
		return ret, map[string]interface{}{}, errors.New("Could not convert")
	}

	return ret, retBody, nil
}
