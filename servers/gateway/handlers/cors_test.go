package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServeHTTP(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	rh := NewResponseHeader(handler)
	cases := []struct {
		name               string
		query              string
		expectedStatusCode int
	}{
		{
			"Valid Method",
			"POST",
			http.StatusOK,
		},
		{
			"Invalid Method",
			"OPTIONS",
			http.StatusOK,
		},
	}

	for _, c := range cases {
		var buffer bytes.Buffer
		req, err := http.NewRequest(c.query, "/v1/users/", bytes.NewBuffer(buffer.Bytes()))
		if err != nil {
			t.Errorf("request error: %v", err)
		}

		rr := httptest.NewRecorder()
		rh.ServeHTTP(rr, req)
		resp := rr.Result()
		if rr.HeaderMap.Get(headerAccessControlAllowOrigin) != "*" {
			t.Errorf("expected to get * as Access-Control-Allow-Origin, but got %s", rr.HeaderMap.Get(headerAccessControlAllowOrigin))
		}

		if rr.HeaderMap.Get("Access-Control-Allow-Methods") != "GET, PUT, POST, PATCH, DELETE" {
			t.Errorf("expected to get GET, PUT, POST, PATCH, DELETE as Access-Control-Allow-Methods, but got %s", rr.HeaderMap.Get("Access-Control-Allow-Methods"))
		}

		if rr.HeaderMap.Get("Access-Control-Allow-Headers") != (headerContentType + ", " + "Authorization") {
			t.Errorf("expected to get %s as Access-Control-Allow-Headers, but got %s", (headerContentType + ", " + "Authorization"), rr.HeaderMap.Get("Access-Control-Allow-Headers"))
		}

		if rr.HeaderMap.Get("Access-Control-Expose-Headers") != "Authorization" {
			t.Errorf("expected to get Authorization as Access-Control-Expose-Headers, but got %s", rr.HeaderMap.Get("Access-Control-Expose-Headers"))
		}

		if rr.HeaderMap.Get("Access-Control-Max-Age") != "600" {
			t.Errorf("expected to get 600 as Access-Control-Max-Age, but got %s", rr.HeaderMap.Get("Access-Control-Max-Age"))
		}

		if resp.StatusCode != c.expectedStatusCode {
			t.Errorf("case %s: incorrect status code: expected %d but got %d", c.name, c.expectedStatusCode, resp.StatusCode)
		}
	}
}