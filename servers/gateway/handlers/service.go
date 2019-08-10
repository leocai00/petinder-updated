package handlers

import (
	"github.com/final-project-petinder/servers/gateway/sessions"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"strings"
	"net/http/httputil"
)

// NewServiceProxy returns a new ReverseProxy
func (ctx *MyHandler) NewServiceProxy(addr string) *httputil.ReverseProxy {
	splitAddr := strings.Split(addr, ",")
	nextAddr := 0
	mx := sync.Mutex{}

	return &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = "http"
			mx.Lock()
			r.URL.Host = splitAddr[nextAddr]
			nextAddr = (nextAddr + 1) % len(splitAddr)
			mx.Unlock()

			ss := &SessionState{}
			
			_, err := sessions.GetState(r, ctx.Key, ctx.SessionStore, ss)
			if err != nil {
				log.Printf(fmt.Sprintf("session id error: %v", err.Error()))
				return
			}

			r.Header.Del("X-User")

			userJSON, err := json.Marshal(ss.Users)
			if err != nil {
				log.Printf(fmt.Sprintf("marshal error: %v", err.Error()))
				return
			}

			log.Printf("user json: %v", string(userJSON))
			r.Header.Add("X-User", string(userJSON))
		},
	}
}