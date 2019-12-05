package handlers

import (
	"INFO441-RoomReservation/servers/gateway/sessions"
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"sync"
)

type UserLite struct {
	ID       int64  `json:"userID"`
	UserName string `json:"userName"`
	Type     string `json:"userType"`
}

// this proxy checks the authentication header
// if the token is authenticated
// pass the user information into reservation system
// if not, remove X-User Header
func (ctx *HandlerContext) NewServiceProxy(addr string) *httputil.ReverseProxy {
	mx := sync.Mutex{}

	return &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = "http"
			mx.Lock()
			r.URL.Host = addr
			mx.Unlock()

			r.Header.Del("X-User")
			userStore := &UserLite{}
			_, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, userStore)
			if err != nil {
				return
			}
			userJSON, err := json.Marshal(userStore)
			if err != nil {
				return
			}

			r.Header.Add("X-User", string(userJSON))
		},
	}
}
