package handlers

import (
	M "INFO441-RoomReservation/servers/reservation/models"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func RoomSearchHandler(w http.ResponseWriter, r *http.Request) {
	user := Authenticate(r)
	if user == nil {
		http.Error(w, "Please Sign in", http.StatusUnauthorized)
		return
	}

	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		log.Printf("Cannot read request body: %v \n", readErr)
	}
	searchTerms := &M.SearchTerm{}
	json.Unmarshal(body, &searchTerms)

}

func RoomReserveHandler(w http.ResponseWriter, r *http.Request) {

}
