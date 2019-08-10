package handlers

import (
	"github.com/cchen97/final-project-petinder/servers/gateway/models/users"
	"github.com/cchen97/final-project-petinder/servers/gateway/sessions"
	"strings"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUsersHandler(t *testing.T) {
	handler := &MyHandler{
		Key: "default",
		SessionStore: sessions.NewMemStore(time.Hour, time.Minute),
		UserStore: &users.FakeStore{},
	}
	cases := []struct {
		name                string
		query               string
		expectedStatusCode  int
		expectedContentType string
		user                *users.NewUser
	}{
		{
			"Valid Method",
			http.MethodPost,
			http.StatusCreated,
			applicationJSON,
			&users.NewUser{
				Email:        "test@example.com",
				Password:     "password",
				PasswordConf: "password",
				UserName:     "test",
				FirstName:    "firstName",
				LastName:     "lastName",
			},
		},
		{
			"Invalid Method",
			http.MethodPatch,
			http.StatusMethodNotAllowed,
			textPlain,
			&users.NewUser{
				Email:        "test@example.com",
				Password:     "password",
				PasswordConf: "password",
				UserName:     "test",
				FirstName:    "firstName",
				LastName:     "lastName",
			},
		},
		{
			"Improper Header Content Type",
			http.MethodPost,
			http.StatusUnsupportedMediaType,
			textPlain,
			&users.NewUser{
				Email:        "test@example.com",
				Password:     "password",
				PasswordConf: "password",
				UserName:     "test",
				FirstName:    "firstName",
				LastName:     "lastName",
			},
		},
		{
			"Invalid New User",
			http.MethodPost,
			http.StatusInternalServerError,
			textPlain,
			&users.NewUser{
				Email:        "test@example.com",
				Password:     "password",
				PasswordConf: "notmatch",
				UserName:     "test",
				FirstName:    "firstName",
				LastName:     "lastName",
			},
		},
	}

	for _, c := range cases {
		body, _ := json.Marshal(c.user)
		req := httptest.NewRequest(c.query, "/v1/users", strings.NewReader(string(body)))

		if c.expectedStatusCode != http.StatusUnsupportedMediaType {
			req.Header.Add(headerContentType, applicationJSON)
		}

		rr := httptest.NewRecorder()
		handler.UsersHandler(rr, req)
		resp := rr.Result()

		if resp.StatusCode != c.expectedStatusCode {
			t.Errorf("case %s: incorrect status code: expected %d but got %d",
				c.name, c.expectedStatusCode, resp.StatusCode)
		}

		ct := resp.Header.Get(headerContentType)
		if ct != c.expectedContentType {
			t.Errorf("case %s: incorrect Content-Type header: expected %s but got %s",
				c.name, c.expectedContentType, ct)
		}
	}
}

func TestSpecificUsersHandler(t *testing.T) {
	state := &SessionState{
		Time: time.Now(),
		Address: "blah",
		Users: fakeUser(),
	}

	handler := &MyHandler{
		Key: "default",
		SessionStore: sessions.NewMemStore(time.Hour, time.Minute),
		UserStore: &users.FakeStore{},
	}

	update := &users.Updates{
		FirstName: "newFirstName",
		LastName:  "newLastName",
	}

	cases := []struct {
		name                string
		query               string
		expectedStatusCode  int
		expectedContentType string
		updates             *users.Updates
		handler             *MyHandler
		auth                bool
		ss                  *SessionState
		id                  string
	}{
		{
			"Invalid URL",
			http.MethodGet,
			http.StatusInternalServerError,
			textPlain,
			update,
			handler,
			false,
			state,
			"blah/wrong",
		},
		{
			"Unauthorized User",
			http.MethodGet,
			http.StatusUnauthorized,
			textPlain,
			update,
			handler,
			false,
			state,
			"10",
		},
		{
			"Invalid Method",
			http.MethodDelete,
			http.StatusMethodNotAllowed,
			textPlain,
			update,
			handler,
			true,
			state,
			"55",
		},
		{
			"Successful Request with Resource Path 'me'",
			http.MethodGet,
			http.StatusOK,
			applicationJSON,
			update,
			handler,
			true,
			state,
			"me",
		},
		{
			"Improper Header Content Type",
			http.MethodPatch,
			http.StatusUnsupportedMediaType,
			textPlain,
			update,
			handler,
			true,
			state,
			"me",
		},
		{
			"Error Updating",
			http.MethodPatch,
			http.StatusInternalServerError,
			textPlain,
			&users.Updates{
				FirstName: "",
				LastName: "",
			},
			handler,
			true,
			state,
			"me",
		},
	}

	for _, c := range cases {
		var buffer bytes.Buffer

		err := json.NewEncoder(&buffer).Encode(c.updates)
		if err != nil {
			t.Errorf("encode error: %v", err)
		}

		req, err := http.NewRequest(c.query, "/v1/users/" + c.id, bytes.NewBuffer(buffer.Bytes()))
		if err != nil {
			t.Errorf("request error: %v", err)
		}

		if c.expectedStatusCode != http.StatusUnsupportedMediaType {
			req.Header.Add(headerContentType, applicationJSON)
		}

		rr := httptest.NewRecorder()
		sid, _ := sessions.BeginSession(c.handler.Key, c.handler.SessionStore, c.ss, rr)
		if c.auth {
			req.Header.Add("Authorization", "Bearer "+sid.String())
		}

		if c.expectedStatusCode != http.StatusUnsupportedMediaType {
			req.Header.Add(headerContentType, applicationJSON)
		}

		handler.SpecificUserHandler(rr, req)
		resp := rr.Result()

		if resp.StatusCode != c.expectedStatusCode {
			t.Errorf("case %s: incorrect status code: expected %d but got %d",
				c.name, c.expectedStatusCode, resp.StatusCode)
		}

		if resp.Header.Get(headerContentType) != c.expectedContentType {
			t.Errorf("case %s: incorrect header content type: expected %s but got %s",
				c.name, c.expectedContentType, resp.Header.Get(headerContentType))
		}
	}
}

func fakeUser() *users.User {
	user := &users.NewUser{
		Email:        "fake@gmail.com",
		Password:     "fakepassword",
		PasswordConf: "fakepassword",
		UserName:     "newfake",
		FirstName:    "fakeFirstName",
		LastName:     "fakeLastName",
	}
	u, _ := user.ToUser()
	return u
}

func TestSessionHandler(t *testing.T) {
	handler := &MyHandler{
		Key: "default",
		SessionStore: sessions.NewMemStore(time.Hour, time.Minute),
		UserStore: &users.FakeStore{},
	}

	cases := []struct {
		name                string
		query               string
		expectedStatusCode  int
		expectedContentType string
		cred                *users.Credentials
		handler             *MyHandler
	}{
		{
			"Valid Method",
			http.MethodPost,
			http.StatusCreated,
			applicationJSON,
			&users.Credentials{
				Email:    "fake@gmail.com",
				Password: "fakepassword",
			},
			handler,
		},
	}

	for _, c := range cases {
		var buffer bytes.Buffer

		err := json.NewEncoder(&buffer).Encode(c.cred)
		if err != nil {
			t.Errorf("encode error: %v", err)
		}

		req, err := http.NewRequest(c.query, "/v1/sessions", bytes.NewBuffer(buffer.Bytes()))
		if err != nil {
			t.Errorf("request error: %v", err)
		}

		if c.expectedStatusCode != http.StatusUnsupportedMediaType {
			req.Header.Add(headerContentType, applicationJSON)
		}

		rr := httptest.NewRecorder()
		c.handler.SessionsHandler(rr, req)
		resp := rr.Result()

		if resp.StatusCode != c.expectedStatusCode {
			t.Errorf("case %s: incorrect status code: expected %d but got %d",
				c.name, c.expectedStatusCode, resp.StatusCode)
		}
		
		if resp.Header.Get(headerContentType) != c.expectedContentType {
			t.Errorf("case %s: incorrect header content type: expected %s but got %s",
				c.name, c.expectedContentType, resp.Header.Get(headerContentType))
		}
	}
}

func TestSpecificSessionHandler(t *testing.T) {
	handler := &MyHandler{
		Key: "default",
		SessionStore: sessions.NewMemStore(time.Hour, time.Minute),
		UserStore: &users.FakeStore{},
	}
	cases := []struct {
		name                string
		query               string
		expectedStatusCode  int
		expectedContentType string
		q                   string
	}{
		{
			"Valid Method",
			http.MethodDelete,
			0,
			textPlain,
			"mine",
		},
		{
			"Invalid Method",
			http.MethodPatch,
			http.StatusMethodNotAllowed,
			textPlain,
			"mine",
		},
		{
			"Invalid URL or Unauthenticated User",
			http.MethodDelete,
			http.StatusInternalServerError,
			textPlain,
			"mine",
		},
	}

	for _, c := range cases {
		var buffer bytes.Buffer

		req, err := http.NewRequest(c.query, "/v1/sessions/"+c.q, bytes.NewBuffer(buffer.Bytes()))
		if err != nil {
			t.Errorf("request error: %v", err)
		}

		req.Header.Add(headerContentType, textPlain)
		rr := httptest.NewRecorder()
		handler.SpecificSessionHandler(rr, req)
		resp := rr.Result()

		if resp.StatusCode != c.expectedStatusCode && c.name != "Valid Method" {
			t.Errorf("case %s: incorrect status code: expected %d but got %d",
				c.name, c.expectedStatusCode, resp.StatusCode)
		}

		if resp.Header.Get(headerContentType) != c.expectedContentType {
			t.Errorf("case %s: incorrect header content type: expected %s but got %s",
				c.name, c.expectedContentType, resp.Header.Get(headerContentType))
		}
	}
}