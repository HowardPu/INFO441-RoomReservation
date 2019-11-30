package handlers

import (
	"INFO441-RoomReservation/servers/gateway/sessions"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"sync/atomic"
)

type UserLite struct {
	ID       int64  `json:"userID"`
	UserName string `json:"userName"`
	Type     string `json:"userType"`
}

func (ctx *HandlerContext) CustomRouting(targets []string) func(r *http.Request) {
	var counter int32 = 0
	// mx := sync.Mutex{}
	return func(r *http.Request) {
		length := int32(len(targets))
		targ := targets[counter%length]
		atomic.AddInt32(&counter, 1)

		r.URL.Scheme = "http"
		// mx.Lock()
		r.URL.Host = targ
		// mx.Unlock()

		r.Header.Del("X-User")
		userStore := &UserLite{}
		_, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, userStore)
		if err != nil {
			log.Printf("getState err: %v /n", err)
			return
		}
		userJSON, err := json.Marshal(userStore)
		if err != nil {
			log.Printf("marshal error: %v /n", err)
			return
		}

		r.Header.Add("X-User", string(userJSON))
	}
}

func (ctx *HandlerContext) NewServiceProxy(targets []string) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Director: ctx.CustomRouting(targets),
	}
}
