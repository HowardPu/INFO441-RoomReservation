package sessions

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
)

const headerAuthorization = "Authorization"
const paramAuthorization = "auth"
const schemeBearer = "Bearer "

//ErrNoSessionID is used when no session ID was found in the Authorization header
var ErrNoSessionID = errors.New("no session ID found in " + headerAuthorization + " header")

//ErrInvalidScheme is used when the authorization scheme is not supported
var ErrInvalidScheme = errors.New("authorization scheme not supported")

//BeginSession creates a new SessionID, saves the `sessionState` to the store, adds an
//Authorization header to the response with the SessionID, and returns the new SessionID
func BeginSession(signingKey string, store Store, sessionState interface{}, w http.ResponseWriter) (SessionID, error) {
	//TODO:
	//- create a new SessionID
	//- save the sessionState to the store
	//- add a header to the ResponseWriter that looks like this:
	//    "Authorization: Bearer <sessionID>"
	//  where "<sessionID>" is replaced with the newly-created SessionID
	//  (note the constants declared for you above, which will help you avoid typos)

	newSessionID, err := NewSessionID(signingKey)
	if err != nil {
		return InvalidSessionID, err
	}

	store.Save(newSessionID, sessionState)
	w.Header().Set(headerAuthorization, schemeBearer+newSessionID.String())

	return newSessionID, nil
}

//GetSessionID extracts and validates the SessionID from the request headers
func GetSessionID(r *http.Request, signingKey string) (SessionID, error) {
	//TODO: get the value of the Authorization header,
	//or the "auth" query string parameter if no Authorization header is present,
	//and validate it. If it's valid, return the SessionID. If not
	//return the validation error.

	reqToken := r.Header.Get(headerAuthorization)
	id := ""

	if len(reqToken) == 0 {
		query, _ := r.URL.Query()[paramAuthorization]
		if len(query) == 0 {
			return InvalidSessionID, ErrNoSessionID
		}

		if len(query) > 1 {
			return SessionID(id), errors.New("More than one Session IDs")
		}

		reqToken = query[0]
	}

	reqToken = strings.TrimSpace(reqToken)

	re := regexp.MustCompile(" +")
	tokenSplit := re.Split(reqToken, -1)

	if len(tokenSplit) <= 1 {
		return InvalidSessionID, ErrInvalidScheme
	}

	if len(tokenSplit) > 2 {
		return InvalidSessionID, errors.New("Too Many sessionIDs!")
	}

	scheme := tokenSplit[0]

	if scheme+" " != schemeBearer {
		return InvalidSessionID, ErrInvalidScheme
	}

	id = tokenSplit[1]

	sessionID, err := ValidateID(id, signingKey)

	if err != nil {
		return InvalidSessionID, err
	}

	return sessionID, nil
}

//GetState extracts the SessionID from the request,
//gets the associated state from the provided store into
//the `sessionState` parameter, and returns the SessionID
func GetState(r *http.Request, signingKey string, store Store, sessionState interface{}) (SessionID, error) {
	sessionID, err := GetSessionID(r, signingKey)
	if err != nil {
		return InvalidSessionID, err
	}

	getErr := store.Get(sessionID, &sessionState)
	if getErr != nil {
		return InvalidSessionID, getErr
	}

	return sessionID, nil
}

//EndSession extracts the SessionID from the request,
//and deletes the associated data in the provided store, returning
//the extracted SessionID.
func EndSession(r *http.Request, signingKey string, store Store) (SessionID, error) {
	//TODO: get the SessionID from the request, and delete the
	//data associated with it in the store.

	sessionID, err := GetSessionID(r, signingKey)

	if err != nil {
		return InvalidSessionID, err
	}

	deleteErr := store.Delete(sessionID)
	return sessionID, deleteErr
}
