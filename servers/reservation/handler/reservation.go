package handler

import (
	"net/http"
)

func RoomSearchHandler(w http.ResponseWriter, r *http.Request) {
	user := Authenticate(r)
	if user == nil {
		http.Error(w, "Please Sign in", http.StatusUnauthorized)
		return
	}
}

func RoomReserveHandler(w http.ResponseWriter, r *http.Request) {

}
