package handlers

import (
	M "INFO441-RoomReservation/servers/reservation/models"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/streadway/amqp"
)

type RoomForm struct {
	RoomName     string `json:"roomName"`
	Floor        int    `json:"floor"`
	Capcity      int    `json:"capacity"`
	RoomTypeName string `json:"roomType"`
}

type ReserveForm struct {
	Year      string `json:"year"`
	Month     string `json:"month"`
	Day       string `json:"day"`
	RoomName  string `json:"roomName"`
	BeginTime int    `json:"beginTime"`
	Duration  int    `json:"duration"`
}

var ReservationMethods map[string]bool = map[string]bool{
	http.MethodPost:   true,
	http.MethodDelete: true,
	http.MethodGet:    true,
}

var RoomMethods map[string]bool = map[string]bool{
	http.MethodPost:   true,
	http.MethodPatch:  true,
	http.MethodGet:    true,
	http.MethodDelete: true,
}

var IssueMethods map[string]bool = map[string]bool{
	http.MethodPost:  true,
	http.MethodPatch: true,
	http.MethodGet:   true,
}

var EquipMethods map[string]bool = map[string]bool{
	http.MethodPost:   true,
	http.MethodGet:    true,
	http.MethodDelete: true,
}

func RoomSearchHandler(w http.ResponseWriter, r *http.Request) {
	user := Authenticate(r)
	if user == nil {
		http.Error(w, "Please Sign in", http.StatusUnauthorized)
		return
	}

	headerErr := CheckContentType(&w, r)
	if headerErr != nil {
		return
	}

	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		log.Printf("Cannot read request body: %v \n", readErr)
	}
	searchTerms := &M.SearchTerm{}
	json.Unmarshal(body, &searchTerms)

}

// What I need in the body:
// 1: year
// 2: month
// 3: day
// 4: room name
// 4: begin time
// 5: duration

func (ctx *HandlerContext) RoomReserveHandler(w http.ResponseWriter, r *http.Request) {
	user := Authenticate(r)
	if user == nil {
		http.Error(w, "Please Sign in", http.StatusUnauthorized)
		return
	}

	headerErr := CheckContentType(&w, r)
	if headerErr != nil {
		return
	}

	methodErr := CheckMethods(&w, r, ReservationMethods)
	if methodErr != nil {
		return
	}

	body, readErr := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if readErr != nil {
		http.Error(w, "Cannot Read Request", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodPost {
		newResForm := ReserveForm{}

		marshalErr := json.Unmarshal(body, &newResForm)
		if marshalErr != nil {
			http.Error(w, "Cannot Unmarshal Request", http.StatusBadRequest)
			return
		}

		date := newResForm.Year + "-" + newResForm.Month + "-" + newResForm.Day
		userName := user.Name
		newID, tranErr := ctx.ReservationStore.ReserveRoom(userName, newResForm.RoomName,
			newResForm.BeginTime, newResForm.Duration, date)
		if tranErr != nil {
			message := fmt.Sprintf("Cannot Reserve: %v", tranErr)
			http.Error(w, message, http.StatusBadRequest)
			return
		}

		mqMessage := fmt.Sprintf(`{"type":"reservation-create", "roomName":%v, "begin": %d, "duration": %d}`,
			newResForm.RoomName, newResForm.BeginTime, newResForm.Duration)

		ctx.PublishMessage(mqMessage)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		resultBody := fmt.Sprintf(`{"id": %d}`, newID)
		w.Write([]byte(resultBody))
	} else {
		deleteForm := map[string]int{}
		marshalErr := json.Unmarshal(body, &deleteForm)
		if marshalErr != nil {
			http.Error(w, "Cannot Unmarshal Request", http.StatusBadRequest)
			return
		}
		resID, found := deleteForm["id"]
		if !found {
			http.Error(w, "Reservation ID not found", http.StatusBadRequest)
			return
		}

		userName := user.Name

		roomName, deleteError := ctx.ReservationStore.ReleaseReservation(userName, resID)

		if deleteError != nil {
			message := fmt.Sprintf("Cannot Reserve: %v", deleteError)
			http.Error(w, message, http.StatusBadRequest)
			return
		}
		mqMessage := fmt.Sprintf(`{"type":"reservation-delete", "roomName":%v}`, roomName)
		ctx.PublishMessage(mqMessage)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Release Reservation!"))
	}
}

func (ctx *HandlerContext) RoomHandler(w http.ResponseWriter, r *http.Request) {

	user := Authenticate(r)
	if user == nil {
		http.Error(w, "Please Sign in", http.StatusUnauthorized)
		return
	}

	headerErr := CheckContentType(&w, r)
	if headerErr != nil {
		return
	}

	methodErr := CheckMethods(&w, r, RoomMethods)
	if methodErr != nil {
		return
	}

	body, readErr := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if readErr != nil {
		http.Error(w, "Cannot Read Request", http.StatusBadRequest)
		return
	}

	userName := user.Name

	if r.Method == http.MethodPost {

		newRoomForm := RoomForm{}

		marshalErr := json.Unmarshal(body, &newRoomForm)
		if marshalErr != nil {
			http.Error(w, "Cannot Unmarshal Request", http.StatusBadRequest)
			return
		}

		newID, tranErr := ctx.ReservationStore.AddRoom(newRoomForm.RoomName, newRoomForm.Floor,
			newRoomForm.Capcity, newRoomForm.RoomTypeName, userName)
		if tranErr != nil {
			message := fmt.Sprintf("Cannot Create Room: %v", tranErr)
			http.Error(w, message, http.StatusBadRequest)
			return
		}

		roomData, _ := json.Marshal(newRoomForm)

		mqMessage := fmt.Sprintf(`{"type":"room-create", "room":%v}`, string(roomData))
		ctx.PublishMessage(mqMessage)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		resultBody := fmt.Sprintf(`{"id": %d}`, newID)
		w.Write([]byte(resultBody))
	} else {
		// delete room
	}

}

func (ctx *HandlerContext) PublishMessage(message string) {
	ctx.RabbitConnection.Publish(
		"",                  // exchange
		ctx.RabbitQueueName, // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

func CheckContentType(w *http.ResponseWriter, r *http.Request) error {
	contentHeader := r.Header.Get("Content-type")
	if contentHeader != "application/json" {
		http.Error(*w, "Request needs to have application/json header", http.StatusUnsupportedMediaType)
		return errors.New("")
	}
	return nil
}

func CheckMethods(w *http.ResponseWriter, r *http.Request, allowedMethods map[string]bool) error {
	_, found := allowedMethods[r.Method]
	if !found {
		http.Error(*w, "Wrong Method!", http.StatusMethodNotAllowed)
		return errors.New("")
	}
	return nil
}
