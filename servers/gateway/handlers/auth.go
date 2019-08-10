package handlers

import (
	"github.com/final-project-petinder/servers/gateway/sessions"
	"github.com/final-project-petinder/servers/gateway/models/users"
	"sort"
	"strconv"
	"path"
	"time"
	"fmt"
	"encoding/json"
	"strings"
	"net/http"
)

//TODO: define HTTP handler functions as described in the
//assignment description. Remember to use your handler context
//struct as the receiver on these functions so that you have
//access to things like the session store and user store.
const headerContentType = "Content-Type"
const headerAccessControlAllowOrigin = "Access-Control-Allow-Origin"
const textPlain = "text/plain; charset=utf-8"
const applicationJSON = "application/json"

// UsersHandler handles requests for the "users" resource.
func (ctx *MyHandler) UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if !strings.HasPrefix(r.Header.Get(headerContentType), applicationJSON) {
			http.Error(w, "The request body must be in JSON.", http.StatusUnsupportedMediaType)
			return
		}

		var nu *users.NewUser

		err := json.NewDecoder(r.Body).Decode(&nu)
		if err != nil {
			http.Error(w, "Error decoding.", http.StatusBadRequest)
			return
		}
		
		u, err := nu.ToUser()
		if err != nil {
			http.Error(w, "Error converting user.", http.StatusBadRequest)
			return
		}

		user, err := ctx.UserStore.Insert(u)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error inserting user: %v", err), http.StatusBadRequest)
			return
		}

		ctx.Trie.Add(user.UserName, user.ID)
		ctx.Trie.Add(user.LastName, user.ID)
		ctx.Trie.Add(user.FirstName, user.ID)

		s := &SessionState{time.Now(), r.RemoteAddr, user}
		_, err = sessions.BeginSession(ctx.Key, ctx.SessionStore, s, w)
		if err != nil {
			http.Error(w, "Error beginning session.", http.StatusInternalServerError)
			return
		}

		w.Header().Add(headerContentType, applicationJSON)
		w.WriteHeader(http.StatusCreated)

		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			http.Error(w, "Error encoding.", http.StatusInternalServerError)
			return
		}
	} else if r.Method == "GET" {		
		ss := &SessionState{}
		_, err := sessions.GetState(r, ctx.Key, ctx.SessionStore, ss)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error getting session state: %v", err), http.StatusUnauthorized)
			return
		}

		q := r.FormValue("q")
		if q == "" {
			http.Error(w, "Empty query.", http.StatusBadRequest)
			return
		}

		var user []*users.User
		arr := ctx.Trie.Find(20, q)
		for i := range arr {
			u, err := ctx.UserStore.GetByID(arr[i])
			if err != nil {
				http.Error(w, fmt.Sprintf("Error getting by ID: %v", err), http.StatusInternalServerError)
				return
			}
			user = append(user, u)
		}
		sort.Slice(user, func(i, j int) bool {
			return user[i].UserName < user[j].UserName
		})

		w.Header().Add(headerContentType, applicationJSON)
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error encoding: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}
}

// SpecificUserHandler handles requests for a specific user.
func (ctx *MyHandler) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "PATCH" {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}
	
	url := strings.Split(r.URL.String(), "/")
	if len(url) != 4 || url[1] != "v1" || url[2] != "users" {
		http.Error(w, "Invalid URL.", http.StatusUnauthorized)
		return
	}

	ss := &SessionState{}
	_, err := sessions.GetState(r, ctx.Key, ctx.SessionStore, ss)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting session state: %v", err), http.StatusUnauthorized)
		return
	}

	var uid int64
	if path.Base(r.URL.Path) != "me" {  // given UserID is a Int
		uid, err = strconv.ParseInt(path.Base(r.URL.Path), 10, 64)
		if err != nil {
			http.Error(w, "Error converting ID.", http.StatusUnauthorized)
			return
		}
	} else {  // given UserID is "me"
		uid = ss.Users.ID
	}

	if r.Method == "GET" {
		u, err := ctx.UserStore.GetByID(uid)
		if err != nil {
			http.Error(w, "User not found.", http.StatusNotFound)
			return
		}
		
		w.Header().Add(headerContentType, applicationJSON)
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(u)
		if err != nil {
			http.Error(w, "Error encoding.", http.StatusInternalServerError)
			return
		}
	} else if r.Method == "PATCH" {
		if uid != ss.Users.ID {
			http.Error(w, "User does not match.", http.StatusForbidden)
			return
		}
		
		if !strings.HasPrefix(r.Header.Get(headerContentType), applicationJSON) {
			http.Error(w, "The request body must be in JSON.", http.StatusUnsupportedMediaType)
			return
		}

		u, err := ctx.UserStore.GetByID(uid)
		if err != nil {
			http.Error(w, "User not found.", http.StatusNotFound)
			return
		}

		var updates *users.Updates
		err = json.NewDecoder(r.Body).Decode(&updates)
		if err != nil {
			http.Error(w, "Error decoding.", http.StatusBadRequest)
			return
		}

		oldFirstName := u.FirstName
		oldLastName := u.LastName
		err = u.ApplyUpdates(updates)
		if err != nil {
			http.Error(w, "Invalid updates.", http.StatusBadRequest)
			return
		}

		updatedUser, err := ctx.UserStore.Update(ss.Users.ID, updates)
		if err != nil {
			http.Error(w, "Error updating user.", http.StatusInternalServerError)
			return
		}

		ctx.Trie.Remove(oldFirstName, ss.Users.ID)
		ctx.Trie.Remove(oldLastName, ss.Users.ID)
		ctx.Trie.Add(updatedUser.FirstName, ss.Users.ID)
		ctx.Trie.Add(updatedUser.LastName, ss.Users.ID)

		w.Header().Add(headerContentType, applicationJSON)
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(updatedUser)
		if err != nil {
			http.Error(w, "Error encoding.", http.StatusInternalServerError)
			return
		}
	}
}

// SessionsHandler handles requests for the "sessions" resource, and allows
// clients to begin a new session using an existing user's credentials.
func (ctx *MyHandler) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if !strings.HasPrefix(r.Header.Get(headerContentType), applicationJSON) {
			http.Error(w, "The request body must be in JSON.", http.StatusUnsupportedMediaType)
			return
		}

		credentials := &users.Credentials{}
		err := json.NewDecoder(r.Body).Decode(credentials)
		if err != nil {
			http.Error(w, "Error decoding.", http.StatusInternalServerError)
			return
		}

		u, err := ctx.UserStore.GetByEmail(credentials.Email)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid credentials: %v", err), http.StatusUnauthorized)
			return
		}

		err = u.Authenticate(credentials.Password)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid credentials: %v", err), http.StatusUnauthorized)
			return
		}

		var addr string
		if r.Header.Get("X-Forwarded-For") != "" {
			addr = r.Header.Get("X-Forwarded-For")
		} else {
			addr = r.RemoteAddr
		}

		ss := &SessionState{time.Now(), addr, u}
		_, err = sessions.BeginSession(ctx.Key, ctx.SessionStore, ss, w)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error beginning session: %v", err), http.StatusInternalServerError)
		}

		w.Header().Add(headerContentType, applicationJSON)
		w.WriteHeader(http.StatusCreated)
		
		err = json.NewEncoder(w).Encode(u)
		if err != nil {
			http.Error(w, "Error encoding.", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}
}

// SpecificSessionHandler handles requests related to a specific authenticated session.
func (ctx *MyHandler) SpecificSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		url := strings.Split(r.URL.String(), "/")
		if len(url) != 4 || url[1] != "v1" || url[2] != "sessions" {
			http.Error(w, "Invalid URL.", http.StatusInternalServerError)
			return
		}

		if path.Base(r.URL.Path) != "mine" {
			http.Error(w, "Unauthenticated user.", http.StatusForbidden)
			return
		}

		_, err := sessions.EndSession(r, ctx.Key, ctx.SessionStore)
		if err != nil {
			http.Error(w, "Error ending session.", http.StatusInternalServerError)
			return
		}

		w.Header().Add(headerContentType, textPlain)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("signed out"))
	} else {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}
}