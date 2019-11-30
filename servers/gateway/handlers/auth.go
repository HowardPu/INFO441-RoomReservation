package handlers

import (
	U "INFO441-RoomReservation/servers/gateway/models/users"
	S "INFO441-RoomReservation/servers/gateway/sessions"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
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
		http.Error(w, "Can't add the create same acccounts twice!", http.StatusBadRequest)
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

// resource path /v1/users/{userID}
func (ctx *HandlerContext) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {

	// throw StatusMethodNotAllowed if the request method is not GET or PATCH
	if r.Method != http.MethodGet && r.Method != http.MethodPatch {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// get user information from session database
	sessionUser := U.User{}
	locker.Lock()
	defer locker.Unlock()
	_, authErr := S.GetState(r, ctx.SessionKey, ctx.SessionStore, &sessionUser)

	// if error occurs, return StatusUnauthorized
	if authErr != nil {
		http.Error(w, "Auth Failed", http.StatusUnauthorized)
		return
	}

	// get id from the path
	path := r.URL.Path
	idSplit := strings.Split(path, "/")
	id := idSplit[len(idSplit)-1]

	user := U.User{}

	idstr := strconv.FormatInt(sessionUser.ID, 10)
	if id != "me" && idstr != id {
		http.Error(w, "User Not Allowed", http.StatusForbidden)
		return
	}

	targetID := sessionUser.ID

	if r.Method == http.MethodGet {
		newUser, getErr := ctx.UserStore.GetById(targetID)
		if getErr != nil {
			http.Error(w, "User Not Found", http.StatusForbidden)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		user = *newUser
	} else {
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			http.Error(w, "Header must be application/json", http.StatusUnsupportedMediaType)
			return
		}

		decoder := json.NewDecoder(r.Body)
		update := U.Updates{}
		decoder.Decode(&update)

		newUser, _ := ctx.UserStore.Update(targetID, &update)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		user = *newUser
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.Encode(user)
	json, _ := json.Marshal(user)
	w.Write(json)
}

func (ctx *HandlerContext) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Not Post Request", http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "Header must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var credential U.Credentials
	decoder.Decode(&credential)

	email := credential.Email
	failAttempt, _ := ctx.SessionStore.GetEmailFailLogIn(email)
	if failAttempt >= 5 {
		http.Error(w, "You Are Locked Due to More Than 5 Attempts", http.StatusUnauthorized)
		return
	}

	locker.RLock()
	user, getErr := ctx.UserStore.GetByEmail(email)
	locker.RUnlock()

	if getErr != nil {
		time.Sleep(1 * time.Second)
		http.Error(w, "Cannot Authenticate", http.StatusUnauthorized)
		return
	}
	authError := user.Authenticate(credential.Password)

	if authError != nil {
		ctx.SessionStore.IncrementFailCount(email)
		http.Error(w, "Cannot Authenticate", http.StatusUnauthorized)
		return
	}

	ctx.SessionStore.RemoveFailRecord(email)

	S.BeginSession(ctx.SessionKey, ctx.SessionStore, user, w)

	userName := user.UserName
	date := time.Now()

	dateStr := date.Format("01-02-2006 15:04:05")

	ip := r.RemoteAddr

	xForwardIP := r.Header.Get("X-Forwarded-For")

	if len(xForwardIP) > 0 {
		ipSplit := strings.Split(xForwardIP, ",")
		ip = strings.TrimSpace(ipSplit[0])
	}

	locker.Lock()
	defer locker.Unlock()
	ctx.UserStore.AddSignInInfo(userName, dateStr, ip)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	encoder := json.NewEncoder(os.Stdout)
	encoder.Encode(user)
	json, _ := json.Marshal(user)
	w.Write(json)
}

func (ctx *HandlerContext) SpecificSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Not Delete Request", http.StatusMethodNotAllowed)
		return
	}

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
