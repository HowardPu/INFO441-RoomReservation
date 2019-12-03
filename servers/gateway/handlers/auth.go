package handlers

import (
	U "INFO441-RoomReservation/servers/gateway/models/users"
	S "INFO441-RoomReservation/servers/gateway/sessions"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var locker = &sync.RWMutex{}

//TODO: define HTTP handler functions as described in the
//assignment description. Remember to use your handler context
//struct as the receiver on these functions so that you have
//access to things like the session store and user store.

func (ctx *HandlerContext) UsersHandler(w http.ResponseWriter, r *http.Request) {
	// throw StatusMethodNotAllowed if the method is not Post
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// throw StatusUnsupportedMediaType if content type is not application/json
	contentType := r.Header.Get("Content-type")
	if contentType != "application/json" {
		http.Error(w, "Header must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	// get new user content
	decoder := json.NewDecoder(r.Body)
	newUser := U.NewUser{}
	err := decoder.Decode(&newUser)

	// if the new user is illegal, throw StatusMethodNotAllowed
	if err != nil || newUser.Validate() != nil {
		http.Error(w, "Wrong for adding new user", http.StatusBadRequest)
		return
	}

	// get user struct from new user
	user, _ := newUser.ToUser()

	// lock the server
	locker.Lock()
	defer locker.Unlock()

	// insert user into user database
	updatedUser, insertErr := ctx.UserStore.Insert(user)

	if insertErr != nil {
		http.Error(w, fmt.Sprintf("%v", insertErr), http.StatusBadRequest)
		return
	}

	// begin sessions
	signingKey := ctx.SessionKey
	S.BeginSession(signingKey, ctx.SessionStore, updatedUser, w)

	// write user information (no email, no pwd) in request body
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	encoder := json.NewEncoder(os.Stdout)
	encoder.Encode(updatedUser)

	json, _ := json.Marshal(updatedUser)

	w.Write(json)
}

// this handler authenticates a user and post its information into the response
func (ctx *HandlerContext) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	// if the method is not posy, throw MethodNotAllowed
	if r.Method != http.MethodPost {
		http.Error(w, "Not Post Request", http.StatusMethodNotAllowed)
		return
	}

	// if request does not have application/json header, throw unsupport media type
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "Header must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	// decode credential information into Credential
	decoder := json.NewDecoder(r.Body)
	var credential U.Credentials
	decodeErr := decoder.Decode(&credential)

	// if something wrong when decode, throw bad request
	if decodeErr != nil {
		fmt.Printf("%v ", decodeErr)
		http.Error(w, "Cannot decode credential", http.StatusBadRequest)
		return
	}

	// get email, check fail attempts, if more than 5, reject authorization
	email := credential.Email
	failAttempt, _ := ctx.SessionStore.GetEmailFailLogIn(email)
	if failAttempt >= 5 {
		http.Error(w, "You Are Locked Due to More Than 5 Attempts", http.StatusUnauthorized)
		return
	}

	// give read lock, read user information
	locker.RLock()
	user, getErr := ctx.UserStore.GetByEmail(email)
	locker.RUnlock()

	// if cannot get user info, sleep 1 sec, and reject authorization
	if getErr != nil {
		time.Sleep(1 * time.Second)
		http.Error(w, "Cannot Authenticate", http.StatusUnauthorized)
		return
	}
	authError := user.Authenticate(credential.Password)

	// if auth failed, increament fail attempt by 1 and throw unauthroized
	if authError != nil {
		ctx.SessionStore.IncrementFailCount(email)
		http.Error(w, "Cannot Authenticate", http.StatusUnauthorized)
		return
	}

	// if seccuss, remove fail record
	ctx.SessionStore.RemoveFailRecord(email)

	// begin session for this user
	locker.Lock()
	S.BeginSession(ctx.SessionKey, ctx.SessionStore, user, w)
	defer locker.Unlock()

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// encode user info into the response
	encoder := json.NewEncoder(os.Stdout)
	encoder.Encode(user)
	json, _ := json.Marshal(user)
	w.Write(json)
}

// this function ends sessuon for a auth token
func (ctx *HandlerContext) SpecificSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Not Delete Request", http.StatusMethodNotAllowed)
		return
	}

	// take tha auth parameter into thr request header
	// 1: delete token in session store
	// 2: remove websocket connection for this token
	path := r.URL.Path
	split := strings.Split(path, "/")
	last := split[len(split)-1]

	if last != "mine" {
		http.Error(w, "Wrong Ending", http.StatusForbidden)
		return
	}

	locker.Lock()
	defer locker.Unlock()
	_, sessionErr := S.EndSession(r, ctx.SessionKey, ctx.SessionStore)
	if sessionErr == nil {
		ctx.SocketStore.RemoveConnection(GetAuthToken(r))
	}
	w.Write([]byte("signed out"))
}

func GetAuthToken(r *http.Request) string {
	reqToken := r.Header.Get("Authorization")

	if len(reqToken) == 0 {
		query, _ := r.URL.Query()["auth"]

		if len(query) == 0 {
			reqToken = ""
		} else {
			reqToken = query[0]
		}
	}

	return reqToken
}
