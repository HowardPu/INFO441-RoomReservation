package handlers

import (
	S "INFO441-RoomReservation/servers/reservation/store"
	"encoding/json"
	"fmt"
	"net/http"
)

// Authenticate current user signin status
// if the request does not have X-User header, return nil
// otherwise, return the user struct
func Authenticate(r *http.Request) *S.User {
	userInfo := r.Header.Get("X-User")

	if len(userInfo) > 0 {
		fmt.Printf("User Logged in! %v \n", userInfo)
		curUser := S.User{}
		json.Unmarshal([]byte(userInfo), &curUser)
		return &curUser
	}
	fmt.Printf("User haven't login")
	return nil
}
