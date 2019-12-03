package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

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
	http.MethodGet:    false,
}

var RoomMethods map[string]bool = map[string]bool{
	http.MethodPost:   true,
	http.MethodGet:    true,
	http.MethodDelete: true,
}

var SpecificRoomMethods map[string]bool = map[string]bool{
	http.MethodPost:   true,
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
	http.MethodGet:    false,
	http.MethodDelete: true,
	http.MethodPatch:  true,
}

func (ctx *HandlerContext) RoomHandler(w http.ResponseWriter, r *http.Request) {
	user := Authenticate(r)
	if user == nil {
		http.Error(w, "Please Sign in", http.StatusUnauthorized)
		return
	}

	checkErr := CheckRequest(&w, r, RoomMethods)

	if checkErr != nil {
		return
	}

	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		http.Error(w, "Cannot read Request Body", http.StatusBadRequest)
		return
	}

	bodyJSON := make(map[string]string)

	marshalErr := json.Unmarshal(body, &bodyJSON)

	if marshalErr != nil {
		http.Error(w, "Cannot marshal Request Body", http.StatusBadRequest)
		return
	}

	roomName, nameFound := bodyJSON["roomName"]

	if !nameFound {
		http.Error(w, "No Room Name!", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodGet || r.Method == http.MethodPost {
		var capacity *int = nil
		var floor *int = nil
		capInReq, capFound := bodyJSON["capacity"]
		if capFound {
			capInt, parseCapErr := strconv.Atoi(capInReq)
			if parseCapErr != nil {
				http.Error(w, "Parse Capacity Field Failed", http.StatusBadRequest)
				return
			}
			capacity = &capInt
		}

		floorInReq, flrFound := bodyJSON["floor"]
		if flrFound {
			flrInt, parseFlrErr := strconv.Atoi(floorInReq)
			if parseFlrErr != nil {
				http.Error(w, "Parse Floor Field Failed", http.StatusBadRequest)
				return
			}
			floor = &flrInt
		}

		var roomType *string = nil

		roomTypeInReq, typeFound := bodyJSON["roomType"]
		if typeFound {
			roomType = &roomTypeInReq
		}

		if r.Method == http.MethodPost {
			if capacity == nil || floor == nil || roomType == nil {
				http.Error(w, "capacity/floor/roomType not in the request", http.StatusBadRequest)
				return
			}

			roomID, insertErr := ctx.ReservationStore.AddRoom(roomName, *floor, *capacity, *roomType, user.Name)
			if insertErr != nil {
				message := fmt.Sprintf("Cannot add Room: %v", insertErr)
				http.Error(w, message, http.StatusBadRequest)
				return
			}
			mqMessage := fmt.Sprintf(`{"type":"room-create", "roomName":%v}`, roomName)
			ctx.PublishMessage(mqMessage)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			resultBody := fmt.Sprintf(`{"id": %d}`, roomID)
			w.Write([]byte(resultBody))
		} else {
			if roomType == nil {
				allType := "*"
				roomType = &allType
			}
			rooms, getErr := ctx.ReservationStore.GetRoomLists(roomName, capacity, floor, *roomType)
			if getErr != nil {
				message := fmt.Sprintf("Cannot search Room: %v", getErr)
				http.Error(w, message, http.StatusBadRequest)
				return
			}

			responseData, marshalErr := json.Marshal(rooms)

			if marshalErr != nil {
				message := fmt.Sprintf("Cannot marshal rooms: %v", marshalErr)
				http.Error(w, message, http.StatusBadRequest)
				return
			}

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(responseData)
		}
		return
	}

	if r.Method == http.MethodDelete {
		roomID, tranErr := ctx.ReservationStore.DeleteRoom(roomName, user.Name)
		if tranErr != nil {
			message := fmt.Sprintf("Cannot Delete rooms: %v", tranErr)
			http.Error(w, message, http.StatusBadRequest)
			return
		}

		mqMessage := fmt.Sprintf(`{"type":"room-delete", "roomID":%d}`, roomID)
		ctx.PublishMessage(mqMessage)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Room is deleted"))
	}
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

	checkErr := CheckRequest(&w, r, ReservationMethods)
	if checkErr != nil {
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
	} else if r.Method == http.MethodDelete {
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
	} else {
		resList, getErr := ctx.ReservationStore.GetReservationLists(user.Name)
		if getErr != nil {
			message := fmt.Sprintf("Cannot get reservation: %v", getErr)
			http.Error(w, message, http.StatusBadRequest)
			return
		}

		responseData, marshalErr := json.Marshal(resList)

		if marshalErr != nil {
			message := fmt.Sprintf("Cannot marshal reservations: %v", marshalErr)
			http.Error(w, message, http.StatusBadRequest)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseData)
	}
}

func (ctx *HandlerContext) EquipmentHandler(w http.ResponseWriter, r *http.Request) {
	user := Authenticate(r)
	if user == nil {
		http.Error(w, "Please Sign in", http.StatusUnauthorized)
		return
	}

	reqErr := CheckRequest(&w, r, EquipMethods)
	if reqErr != nil {
		return
	}

	if r.Method == http.MethodGet {
		result, searchErr := ctx.ReservationStore.GetAllEquipment()
		if searchErr != nil {
			message := fmt.Sprintf("Cannot search equipment: %v", searchErr)
			http.Error(w, message, http.StatusBadRequest)
			return
		}
		jsonStr := "[\"" + strings.Join(result, "\",\"") + "\"]"
		responseBody := fmt.Sprintf(`{"result": %v}`, jsonStr)

		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(responseBody))
		return
	}

	body, readErr := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if readErr != nil {
		http.Error(w, "Cannot Read Request", http.StatusBadRequest)
		return
	}

	jsonData := make(map[string]string)

	unmarshalErr := json.Unmarshal(body, &jsonData)
	if unmarshalErr != nil {
		http.Error(w, "Cannot Unmarshal Request", http.StatusBadRequest)
		return
	}
	equipName, nameFound := jsonData["equipName"]
	if !nameFound {
		http.Error(w, "Equipname not found", http.StatusBadRequest)
		return
	}
	if r.Method == http.MethodPost {
		_, tranErr := ctx.ReservationStore.AddEquipment(equipName, user.Name)
		if tranErr != nil {
			message := fmt.Sprintf("Cannot Add Equipment: %v", tranErr)
			http.Error(w, message, http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Equipment Add!"))
		return
	}

	if r.Method == http.MethodDelete {
		deleteErr := ctx.ReservationStore.DeleteEquipment(equipName, user.Name)
		if deleteErr != nil {
			message := fmt.Sprintf("Cannot delete Equipment: %v", deleteErr)
			http.Error(w, message, http.StatusBadRequest)
			return
		}
		ctx.PublishMessage(`{"type":"equip-delete"}`)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Equipment Deleted!"))
		return
	}

	newName, newNameFound := jsonData["newName"]
	if !newNameFound {
		http.Error(w, "New Name Needed", http.StatusBadRequest)
		return
	}

	updateErr := ctx.ReservationStore.UpdateEquipName(equipName, newName, user.Name)

	if updateErr != nil {
		message := fmt.Sprintf("Cannot update equipment: %v", updateErr)
		http.Error(w, message, http.StatusBadRequest)
		return
	}
	mqMessage := fmt.Sprintf(`{"type": "equip-update", "equipName": %v, "newName": %v}`, equipName, newName)
	ctx.PublishMessage(mqMessage)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Equipment Updated!"))
}

func (ctx *HandlerContext) SpecificRoomHandler(w http.ResponseWriter, r *http.Request) {
	user := Authenticate(r)
	if user == nil {
		http.Error(w, "Please Sign in", http.StatusUnauthorized)
		return
	}

	reqErr := CheckRequest(&w, r, SpecificRoomMethods)
	if reqErr != nil {
		return
	}

	body, readErr := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if readErr != nil {
		http.Error(w, "Cannot Read Request", http.StatusBadRequest)
		return
	}

	jsonData := make(map[string]string)

	unmarshalErr := json.Unmarshal(body, &jsonData)

	if unmarshalErr != nil {
		http.Error(w, "Cannot Unmarshal Request", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodDelete {
		equipRMID, idFound := jsonData["roomEquipID"]
		if !idFound {
			http.Error(w, "No Room Equipment ID", http.StatusBadRequest)
			return
		}
		idInt, parseErr := strconv.Atoi(equipRMID)
		if parseErr != nil {
			http.Error(w, fmt.Sprintf("Cannot Parse ID: %v", equipRMID), http.StatusBadRequest)
			return
		}
		rmErr := ctx.ReservationStore.DeleteEquipmentInRoom(idInt, user.Name)
		if rmErr != nil {
			http.Error(w, fmt.Sprintf("Cannot Remove Equipment for a room: %v", rmErr), http.StatusBadRequest)
			return
		}
		mqMessage := fmt.Sprintf(`{"type": "roomEquip-delete", "id": %d}`, idInt)
		ctx.PublishMessage(mqMessage)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Equip In a room is deleted"))
	} else {
		rmName, nameFound := jsonData["roomName"]
		if !nameFound {
			http.Error(w, "Room Name Required", http.StatusBadRequest)
			return
		}
		if r.Method == http.MethodGet {
			equips, getErr := ctx.ReservationStore.GetEquipList(rmName)
			if getErr != nil {
				http.Error(w, fmt.Sprintf("Cannot Get Equipment for a room: %v", getErr), http.StatusBadRequest)
				return
			}

			responseData, marshalErr := json.Marshal(equips)

			if marshalErr != nil {
				message := fmt.Sprintf("Cannot marshal equips in a room: %v", marshalErr)
				http.Error(w, message, http.StatusBadRequest)
				return
			}

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(responseData)
		} else {
			equipName, equipFound := jsonData["equipName"]
			if !equipFound {
				http.Error(w, "Equip Name Required", http.StatusBadRequest)
				return
			}
			newID, newErr := ctx.ReservationStore.AddEquipmentToRoom(equipName, rmName, user.Name)
			if newErr != nil {
				message := fmt.Sprintf("Cannot add equips in a room: %v", newErr)
				http.Error(w, message, http.StatusBadRequest)
				return
			}
			mqMessage := fmt.Sprintf(`{"type": "roomEquip-add", "roomName":%v, "id": %d}`, rmName, newID)
			ctx.PublishMessage(mqMessage)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("%v is added to %v", equipName, rmName)))
		}
	}
}

func (ctx *HandlerContext) IssueHandler(w http.ResponseWriter, r *http.Request) {
	user := Authenticate(r)
	if user == nil {
		http.Error(w, "Please Sign in", http.StatusUnauthorized)
		return
	}

	reqErr := CheckRequest(&w, r, IssueMethods)
	if reqErr != nil {
		return
	}

	body, readErr := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if readErr != nil {
		http.Error(w, "Cannot Read Request", http.StatusBadRequest)
		return
	}

	jsonData := make(map[string]string)

	unmarshalErr := json.Unmarshal(body, &jsonData)

	if unmarshalErr != nil {
		http.Error(w, "Cannot Unmarshal Request", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodPost {
		issueBody, bodyFound := jsonData["body"]
		if !bodyFound {
			http.Error(w, "No Issue Body", http.StatusBadRequest)
			return
		}
		roomName, nameFound := jsonData["roomName"]
		if !nameFound {
			http.Error(w, "No Room Name", http.StatusBadRequest)
			return
		}
		issueID, addError := ctx.ReservationStore.AddIssue(issueBody, roomName)
		if addError != nil {
			message := fmt.Sprintf("Cannot add issue: %v", addError)
			http.Error(w, message, http.StatusBadRequest)
			return
		}
		mqMessage := fmt.Sprintf(`{"type": "issue-add", "roomName": %v, "issueID": %d}`, roomName, issueID)
		ctx.PublishMessage(mqMessage)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Issue Added"))
	} else {
		typeParam, typeFound := jsonData["type"]
		if !typeFound {
			http.Error(w, "No Type", http.StatusBadRequest)
			return
		}

		roomName, roomNameFound := jsonData["roomName"]
		if !roomNameFound && (r.Method != http.MethodGet || typeParam != "All") {
			http.Error(w, "Room Name Required", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodGet {
			issues, getErr := ctx.ReservationStore.GetIssues(roomName, typeParam)
			if getErr != nil {
				message := fmt.Sprintf("Cannot get issues: %v", getErr)
				http.Error(w, message, http.StatusBadRequest)
				return
			}

			responseData, marshalErr := json.Marshal(issues)

			if marshalErr != nil {
				message := fmt.Sprintf("Cannot marshal issues: %v", marshalErr)
				http.Error(w, message, http.StatusBadRequest)
				return
			}

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(responseData)
		} else {
			issueID, idFound := jsonData["id"]
			if !idFound {
				http.Error(w, "Issue ID required", http.StatusBadRequest)
				return
			}
			idInt, parseError := strconv.Atoi(issueID)
			if parseError != nil {
				http.Error(w, fmt.Sprintf("Cannot parse id: %v", issueID), http.StatusBadRequest)
				return
			}
			updateErr := ctx.ReservationStore.UpdateIssue(idInt, typeParam, user.Name)
			if updateErr != nil {
				http.Error(w, fmt.Sprintf("Cannot update issue: %v", updateErr), http.StatusBadRequest)
				return
			}
			mqMessage := fmt.Sprintf(`{"type": "issue-update", "id": %d, "type": %v}`, idInt, typeParam)
			ctx.PublishMessage(mqMessage)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Issue Updated"))
		}
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

func CheckRequest(w *http.ResponseWriter, r *http.Request, allowedMethods map[string]bool) error {
	needCheck, found := allowedMethods[r.Method]
	if !found {
		http.Error(*w, "Wrong Method!", http.StatusMethodNotAllowed)
		return errors.New("")
	}
	if needCheck {
		contentHeader := r.Header.Get("Content-type")
		if contentHeader != "application/json" {
			http.Error(*w, "Request needs to have application/json header", http.StatusUnsupportedMediaType)
			return errors.New("")
		}
	}
	return nil
}
