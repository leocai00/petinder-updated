package handlers

import (
	"github.com/final-project-petinder/servers/gateway/indexes"
	"github.com/final-project-petinder/servers/gateway/models/users"
	"github.com/final-project-petinder/servers/gateway/sessions"
)

// MyHandler receives HTTP functions
//TODO: define a handler context struct that
//will be a receiver on any of your HTTP
//handler functions that need access to
//globals, such as the key used for signing
//and verifying SessionIDs, the session store
//and the user store
type MyHandler struct {
	Key          string
	SessionStore sessions.Store
	UserStore    users.Store
	Trie         *indexes.Trie
	SocketStore  *SocketStore
}