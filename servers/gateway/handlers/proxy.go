package handlers

import (
	"INFO441-RoomReservation/servers/gateway/sessions"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"sync"
)

type UserLite struct {
	ID       int64  `json:"userID"`
	UserName string `json:"userName"`
	Type     string `json:"userType"`
}

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
				log.Printf("%v /n", err)
				return
			}
			userJSON, err := json.Marshal(userStore)
			if err != nil {
				log.Printf("%v /n", err)
				return
			}

			r.Header.Add("X-User", string(userJSON))
		},
	}
}
