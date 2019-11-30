package handler

import (
	M "INFO441-RoomReservation/servers/reservation/models"
	"encoding/json"
	"fmt"
	"net/http"
)

// Authenticate current user signin status
func Authenticate(r *http.Request) *M.User {
	userInfo := r.Header.Get("X-User")

	if len(userInfo) > 0 {
		fmt.Printf("User Logged in! %v \n", userInfo)
		curUser := &M.User{}
		json.Unmarshal([]byte(userInfo), &curUser)
		return curUser
	}
	fmt.Printf("User haven't login")
	return nil
}
