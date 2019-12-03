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

// for any such may
// key is the allowed method name
// value (boolean) means whether the method required to check application/json header
var ReservationMethods map[string]bool = map[string]bool{
	http.MethodPost:   true,
	http.MethodDelete: true,
	http.MethodGet:    false,
}

var RoomMethods map[string]bool = map[string]bool{
	http.MethodPost:   true,
	http.MethodGet:    false,
	http.MethodDelete: true,
}

var SpecificRoomMethods map[string]bool = map[string]bool{
	http.MethodPost:   true,
	http.MethodGet:    false,
	http.MethodDelete: true,
}

var IssueMethods map[string]bool = map[string]bool{
	http.MethodPost:  true,
	http.MethodPatch: true,
	http.MethodGet:   false,
}

var EquipMethods map[string]bool = map[string]bool{
	http.MethodPost:   true,
	http.MethodGet:    false,
	http.MethodDelete: true,
	http.MethodPatch:  true,
}

var UsedTimeMethod map[string]bool = map[string]bool{
	http.MethodGet: false,
}

/*
	JSON format:
	GET: {
		"roomName": string, ("*" means any)
		"capacity": int (optional)
		"floor": int (optional)
		"roomType": string ("*" means any)
	}

	POST: {
		"roomName": string
		"capacity": int
		"floor": int
		"roomType": string
	}

	DELETE: {
		"roomName": string
	}
*/

func (ctx *HandlerContext) RoomHandler(w http.ResponseWriter, r *http.Request) {
	// if the user is not authenticated, throw unauthorized
	user := Authenticate(r)
	if user == nil {
		http.Error(w, "Please Sign in", http.StatusUnauthorized)
		return
	}

	// check headers, and throw any error
	checkErr := CheckRequest(&w, r, RoomMethods)

	if checkErr != nil {
		return
	}

	// read the request body
	// throw any error
	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		http.Error(w, "Cannot read Request Body", http.StatusBadRequest)
		return
	}

	// marshal the data into map
	// and throw any error
	bodyJSON := make(map[string]string)
	marshalErr := json.Unmarshal(body, &bodyJSON)
	if marshalErr != nil {
		http.Error(w, "Cannot marshal Request Body", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodGet {
		roomType := "*"
		roomName := "*"
		var capacity *int = nil
		var floor *int = nil

		params := r.URL.Query()

		rmTypeInQ, typeFound := params["roomtype"]
		if typeFound {
			if len(rmTypeInQ) != 1 {
				http.Error(w, "Room Type must be one", http.StatusBadRequest)
				return
			}
			roomType = rmTypeInQ[0]
		}

		rmNameInQ, nameFound := params["roomname"]
		if nameFound {
			if len(rmNameInQ) != 1 {
				http.Error(w, "Room Name must be one", http.StatusBadRequest)
				return
			}
			roomName = rmNameInQ[0]
		}

		capInQ, capFound := params["capacity"]
		if capFound {
			if len(capInQ) != 1 {
				http.Error(w, "Capacity must be one", http.StatusBadRequest)
				return
			}
			capStr := capInQ[0]

			capInt, parseErr := strconv.Atoi(capStr)
			if parseErr != nil {
				http.Error(w, "Cannot parse Capacity", http.StatusBadRequest)
				return
			}
			capacity = &capInt
		}

		flrInQ, flrFound := params["floor"]
		if flrFound {
			if len(flrInQ) != 1 {
				http.Error(w, "FLoor must be one", http.StatusBadRequest)
				return
			}
			flrStr := flrInQ[0]

			flrInt, flrErr := strconv.Atoi(flrStr)
			if flrErr != nil {
				http.Error(w, "Cannot parse Floor", http.StatusBadRequest)
				return
			}
			floor = &flrInt
		}

		// Get list of Rooms
		// and throw any error if occured
		rooms, getErr := ctx.ReservationStore.GetRoomLists(roomName, capacity, floor, roomType)
		if getErr != nil {
			message := fmt.Sprintf("Cannot search Room: %v", getErr)
			http.Error(w, message, http.StatusBadRequest)
			return
		}

		// marshal the json data into response
		// throw any error if occured
		responseData, marshalErr := json.Marshal(rooms)

		if marshalErr != nil {
			message := fmt.Sprintf("Cannot marshal rooms: %v", marshalErr)
			http.Error(w, message, http.StatusBadRequest)
			return
		}

		// write application/json header
		// with status ok
		// and write the response (list of target rooms)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseData)
		return
	}

	// get the room name field(requied for all methods besides get)
	// throw bad request if not found
	roomName, nameFound := bodyJSON["roomName"]
	if !nameFound {
		http.Error(w, "No Room Name!", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodPost {
		capInReq, capFound := bodyJSON["capacity"]
		if !capFound {
			http.Error(w, "Capacity required for adding room", http.StatusBadRequest)
			return
		}

		capInt, parseCapErr := strconv.Atoi(capInReq)
		if parseCapErr != nil {
			http.Error(w, "Parse Capacity Field Failed", http.StatusBadRequest)
			return
		}

		floorInReq, flrFound := bodyJSON["floor"]
		// get the floor field in json
		// if exist, parse it and throw any error
		if !flrFound {
			http.Error(w, "Floor required for adding room", http.StatusBadRequest)
			return
		}

		flrInt, parseFlrErr := strconv.Atoi(floorInReq)
		if parseFlrErr != nil {
			http.Error(w, "Parse Floor Field Failed", http.StatusBadRequest)
			return
		}

		roomTypeInReq, typeFound := bodyJSON["roomType"]
		if !typeFound {
			http.Error(w, "Room Type Requied for adding room", http.StatusBadRequest)
			return
		}

		// add new room into the db
		// and throw any error if existed
		roomID, insertErr := ctx.ReservationStore.AddRoom(roomName, flrInt, capInt, roomTypeInReq, user.Name)

		if insertErr != nil {
			message := fmt.Sprintf("Cannot add Room: %v", insertErr)
			http.Error(w, message, http.StatusBadRequest)
			return
		}
		// write room-created in rabbitMq
		mqMessage := fmt.Sprintf(`{"type":"room-create", "roomName":%v, "id":%d}`, roomName, roomID)
		ctx.PublishMessage(mqMessage)

		// write application/json header with status code CREATED
		// write room info in the request body
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		resultBody := fmt.Sprintf(`{"id": %d}`, roomID)
		w.Write([]byte(resultBody))
		return
	}

	if r.Method == http.MethodDelete {
		// delete the room with given room name
		// and throw any error if occured
		roomID, tranErr := ctx.ReservationStore.DeleteRoom(roomName, user.Name)
		if tranErr != nil {
			message := fmt.Sprintf("Cannot Delete rooms: %v", tranErr)
			http.Error(w, message, http.StatusBadRequest)
			return
		}

		// write message to indicate that the room is deleted
		mqMessage := fmt.Sprintf(`{"type":"room-delete", "roomID":%d}`, roomID)
		ctx.PublishMessage(mqMessage)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Room is deleted"))
	}
}

/*

	POST: {
		year: int
		month: int
		day: int
		roomName: string
		beginTime: int
		duration: int
	}

	DELETE: {
		id: int
	}
	GET: nothing
*/
func (ctx *HandlerContext) RoomReserveHandler(w http.ResponseWriter, r *http.Request) {
	// authenticate any user,
	// throw unauthorized if the user not sign in
	user := Authenticate(r)
	if user == nil {
		http.Error(w, "Please Sign in", http.StatusUnauthorized)
		return
	}

	// check request is valid
	// throw any error if occurs
	checkErr := CheckRequest(&w, r, ReservationMethods)
	if checkErr != nil {
		return
	}

	if r.Method == http.MethodGet {
		// get a user's reservation list
		// and return any error
		resList, getErr := ctx.ReservationStore.GetReservationLists(user.Name)
		if getErr != nil {
			message := fmt.Sprintf("Cannot get reservation: %v", getErr)
			http.Error(w, message, http.StatusBadRequest)
			return
		}

		// marshal the data into response body
		// and throw any error
		responseData, marshalErr := json.Marshal(resList)
		if marshalErr != nil {
			message := fmt.Sprintf("Cannot marshal reservations: %v", marshalErr)
			http.Error(w, message, http.StatusBadRequest)
			return
		}

		// set header, status code ok
		// and write the response body
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseData)
	} else {
		// if the method is patch or delete

		// read the body
		// and throw any error if occured
		body, readErr := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if readErr != nil {
			http.Error(w, "Cannot Read Request", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodPost {

			// marshal the json into reservation form
			// and throw any error if occurs
			newResForm := ReserveForm{}
			marshalErr := json.Unmarshal(body, &newResForm)
			if marshalErr != nil {
				http.Error(w, "Cannot Unmarshal Request", http.StatusBadRequest)
				return
			}

			// create date string
			date := newResForm.Year + "-" + newResForm.Month + "-" + newResForm.Day
			userName := user.Name

			// reserve the room at given time spot
			// throw any error if occurs
			newID, tranErr := ctx.ReservationStore.ReserveRoom(userName, newResForm.RoomName,
				newResForm.BeginTime, newResForm.Duration, date)
			if tranErr != nil {
				message := fmt.Sprintf("Cannot Reserve: %v", tranErr)
				http.Error(w, message, http.StatusBadRequest)
				return
			}

			// write the reservation-create message to rabbitmq
			mqMessage := fmt.Sprintf(`{"type":"reservation-create", "roomName":%v, "begin": %d, "duration": %d}`,
				newResForm.RoomName, newResForm.BeginTime, newResForm.Duration)
			ctx.PublishMessage(mqMessage)

			// write the reservation id into response
			// write headers, status created
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			resultBody := fmt.Sprintf(`{"id": %d}`, newID)
			w.Write([]byte(resultBody))
		} else if r.Method == http.MethodDelete {
			// get the id field from json
			// if not exist, throw error
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

			// get username
			// release the reservation
			// and throw any error if occurs
			userName := user.Name
			roomName, deleteError := ctx.ReservationStore.ReleaseReservation(userName, resID)
			if deleteError != nil {
				message := fmt.Sprintf("Cannot Reserve: %v", deleteError)
				http.Error(w, message, http.StatusBadRequest)
				return
			}

			// write "reservation-deleted" in rebbitmq
			// write statusok in response
			mqMessage := fmt.Sprintf(`{"type":"reservation-delete", "roomName":%v}`, roomName)
			ctx.PublishMessage(mqMessage)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Release Reservation!"))
		}
	}
}

/*
	GET: nothing,
	POST: {
		equipName: string
	},
	DELETE: {
		equipName: string
	},
	PATCH: {
		equipName: string
		newName: string
	}

*/
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
		// get all equipments
		// throw any error if occurs
		result, searchErr := ctx.ReservationStore.GetAllEquipment()
		if searchErr != nil {
			message := fmt.Sprintf("Cannot search equipment: %v", searchErr)
			http.Error(w, message, http.StatusBadRequest)
			return
		}

		// encode the result as json into response
		jsonStr := "[\"" + strings.Join(result, "\",\"") + "\"]"
		responseBody := fmt.Sprintf(`{"result": %v}`, jsonStr)

		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(responseBody))
		return
	}

	// if the method is post/patch/delete
	// read body, and throw any error if occurs
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

	// get equipment (required for all methods)
	// throw error if not found
	equipName, nameFound := jsonData["equipName"]
	if !nameFound {
		http.Error(w, "Equipname not found", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodPost {
		// add equipment into db
		// and throw any error if occurs
		_, tranErr := ctx.ReservationStore.AddEquipment(equipName, user.Name)
		if tranErr != nil {
			message := fmt.Sprintf("Cannot Add Equipment: %v", tranErr)
			http.Error(w, message, http.StatusBadRequest)
			return
		}

		// write status ok if success
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Equipment Add!"))
		return
	}

	if r.Method == http.MethodDelete {
		// delete the equipment in db
		// and throw any error if occurs
		deleteErr := ctx.ReservationStore.DeleteEquipment(equipName, user.Name)
		if deleteErr != nil {
			message := fmt.Sprintf("Cannot delete Equipment: %v", deleteErr)
			http.Error(w, message, http.StatusBadRequest)
			return
		}

		// write equip-delete in rabbitmq
		// with status ok
		ctx.PublishMessage(`{"type":"equip-delete"}`)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Equipment Deleted!"))
		return
	}

	// for PATCH
	// find new name field
	// throw error if not found
	newName, newNameFound := jsonData["newName"]
	if !newNameFound {
		http.Error(w, "New Name Needed", http.StatusBadRequest)
		return
	}

	// update the equipment name
	// throw any error occurs
	updateErr := ctx.ReservationStore.UpdateEquipName(equipName, newName, user.Name)
	if updateErr != nil {
		message := fmt.Sprintf("Cannot update equipment: %v", updateErr)
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	// write equip-update in rabbitmq
	// with status ok
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

	// if the method is get
	if r.Method == http.MethodGet {
		// get roomname query
		// throw error if not found
		rmNames, nameFound := r.URL.Query()["roomname"]
		if !nameFound {
			http.Error(w, fmt.Sprintf("Need Room Name for searching equips"), http.StatusBadRequest)
			return
		}

		// if roomname query has more than 1 value
		// throw error
		if len(rmNames) != 1 {
			http.Error(w, fmt.Sprintf("Need only one room name"), http.StatusBadRequest)
			return
		}

		rmName := rmNames[0]

		// get equip informations of a room,
		// and throw any error if occurs
		equips, getErr := ctx.ReservationStore.GetEquipList(rmName)
		if getErr != nil {
			http.Error(w, fmt.Sprintf("Cannot Get Equipment for a room: %v", getErr), http.StatusBadRequest)
			return
		}

		// marshal the equip informations
		// and throw any error if occurs
		responseData, marshalErr := json.Marshal(equips)
		if marshalErr != nil {
			message := fmt.Sprintf("Cannot marshal equips in a room: %v", marshalErr)
			http.Error(w, message, http.StatusBadRequest)
			return
		}

		// write application/json header
		// with status ok
		// and write the response body
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseData)
		return
	}

	// if the method is delete/post
	// read body, and throw any errors if occurs

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
		// for deleting equipment in room
		// search its id
		// throw error if not found/cannot parse
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

		// remove the equipment in a room
		// and throw any error if occurs
		rmErr := ctx.ReservationStore.DeleteEquipmentInRoom(idInt, user.Name)
		if rmErr != nil {
			http.Error(w, fmt.Sprintf("Cannot Remove Equipment for a room: %v", rmErr), http.StatusBadRequest)
			return
		}

		// write roomEquip-delete in th rabbitMQ
		// with status ok
		mqMessage := fmt.Sprintf(`{"type": "roomEquip-delete", "id": %d}`, idInt)
		ctx.PublishMessage(mqMessage)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Equip In a room is deleted"))
	} else {
		// for add equipment to a room
		// search roomName and equipName fields
		// throw anyerror if not found
		rmName, nameFound := jsonData["roomName"]
		if !nameFound {
			http.Error(w, "Room Name Required", http.StatusBadRequest)
			return
		}
		equipName, equipFound := jsonData["equipName"]
		if !equipFound {
			http.Error(w, "Equip Name Required", http.StatusBadRequest)
			return
		}

		// add equipment to a room
		// throw any error that occurs
		newID, newErr := ctx.ReservationStore.AddEquipmentToRoom(equipName, rmName, user.Name)
		if newErr != nil {
			message := fmt.Sprintf("Cannot add equips in a room: %v", newErr)
			http.Error(w, message, http.StatusBadRequest)
			return
		}

		// write roomEquip-add to rabbitMQ
		// with status ok
		mqMessage := fmt.Sprintf(`{"type": "roomEquip-add", "roomName":%v, "id": %d}`, rmName, newID)
		ctx.PublishMessage(mqMessage)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%v is added to %v", equipName, rmName)))
	}
}

func (ctx *HandlerContext) IssueHandler(w http.ResponseWriter, r *http.Request) {

	// get authenticated user
	// if not exist, throw unauthroized
	user := Authenticate(r)
	if user == nil {
		http.Error(w, "Please Sign in", http.StatusUnauthorized)
		return
	}

	// check request, and return any error if occurs
	reqErr := CheckRequest(&w, r, IssueMethods)
	if reqErr != nil {
		return
	}

	if r.Method == http.MethodGet {
		// find roomname/type parameter
		// return any error if not find or too much values
		roomNames, namesFound := r.URL.Query()["roomname"]
		if !namesFound {
			http.Error(w, "Room Name Required", http.StatusBadRequest)
			return
		}
		if len(roomNames) != 1 {
			http.Error(w, "Room Name Required Only one", http.StatusBadRequest)
			return
		}

		roomName := roomNames[0]

		types, typesFound := r.URL.Query()["type"]
		if !typesFound {
			http.Error(w, "Type Params Required", http.StatusBadRequest)
			return
		}
		if len(types) != 1 {
			http.Error(w, "Type Params Required only one", http.StatusBadRequest)
			return
		}
		typeParam := types[0]

		// get issues based on the search type
		// return any error if occurs
		issues, getErr := ctx.ReservationStore.GetIssues(roomName, typeParam)
		if getErr != nil {
			message := fmt.Sprintf("Cannot get issues: %v", getErr)
			http.Error(w, message, http.StatusBadRequest)
			return
		}

		// marshal the data into response
		// with status code ok
		// return any error if occurs
		responseData, marshalErr := json.Marshal(issues)

		if marshalErr != nil {
			message := fmt.Sprintf("Cannot marshal issues: %v", marshalErr)
			http.Error(w, message, http.StatusBadRequest)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseData)
		return
	}

	// for PATCH/Post
	// read request body
	// return any error if occurs
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
		// find issue body
		// if it does not exist/empty, return error
		issueBody, bodyFound := jsonData["body"]
		if !bodyFound || len(issueBody) == 0 {
			http.Error(w, "No Issue Body Or Body is Empty", http.StatusBadRequest)
			return
		}

		// get room name
		// return error if not found
		roomName, nameFound := jsonData["roomName"]
		if !nameFound {
			http.Error(w, "No Room Name", http.StatusBadRequest)
			return
		}

		// add issue to db
		// throw any error if occurs
		// write issue-add message to rabbitmq
		// with status code created
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
		// get updateType (confirm/solve) and issueID
		// throw errors if they are not exist or in bad form
		typeParam, typeFound := jsonData["type"]
		if !typeFound {
			http.Error(w, "No Type", http.StatusBadRequest)
			return
		}

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

		// update issues, throw any error if returned
		updateErr := ctx.ReservationStore.UpdateIssue(idInt, typeParam, user.Name)
		if updateErr != nil {
			http.Error(w, fmt.Sprintf("Cannot update issue: %v", updateErr), http.StatusBadRequest)
			return
		}

		// write status ok if succeed
		// write issue-update to rabbitmq
		mqMessage := fmt.Sprintf(`{"type": "issue-update", "id": %d, "type": %v}`, idInt, typeParam)
		ctx.PublishMessage(mqMessage)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Issue Updated"))

	}
}

func (ctx *HandlerContext) GetUsedTimeHandler(w http.ResponseWriter, r *http.Request) {
	// get authenticated user
	// if not exist, throw unauthroized
	user := Authenticate(r)
	if user == nil {
		http.Error(w, "Please Sign in", http.StatusUnauthorized)
		return
	}

	// check request, and return any error if occurs
	reqErr := CheckRequest(&w, r, UsedTimeMethod)
	if reqErr != nil {
		return
	}

	params := r.URL.Query()

	// get roomname, year, month, day from the query
	// if any of them are not exist or in the bad form
	// throw error
	roomNames, namesFound := params["roomname"]
	if !namesFound {
		http.Error(w, "Room Name Required", http.StatusBadRequest)
		return
	}
	if len(roomNames) != 1 {
		http.Error(w, "Room Name Required Only one", http.StatusBadRequest)
		return
	}
	roomName := roomNames[1]

	years, yearsFound := params["year"]
	if !yearsFound {
		http.Error(w, "Year Required", http.StatusBadRequest)
		return
	}
	if len(years) != 1 {
		http.Error(w, "Year Required Only one", http.StatusBadRequest)
		return
	}
	year := years[0]

	months, monthsFound := params["month"]
	if !monthsFound {
		http.Error(w, "Month Required", http.StatusBadRequest)
		return
	}
	if len(months) != 1 {
		http.Error(w, "Month Required Only one", http.StatusBadRequest)
		return
	}
	month := months[0]

	days, daysFound := params["day"]
	if !daysFound {
		http.Error(w, "Day Required", http.StatusBadRequest)
		return
	}
	if len(days) != 1 {
		http.Error(w, "Day Required Only one", http.StatusBadRequest)
		return
	}
	day := days[0]

	// make date string, search the used time from the db
	// throw any error if occured
	date := year + "-" + month + "-" + day
	result, resultErr := ctx.ReservationStore.GetUsedTime(roomName, date)
	if resultErr != nil {
		message := fmt.Sprintf("Cannot search used time: %v", resultErr)
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	// transform the search result into string
	resStr := "["
	for i := 0; i < len(result); i++ {
		cur := (*result)[i]
		if cur != 0 {
			resStr += strconv.Itoa(i) + ", "
		}
	}
	resStr = resStr + "]"

	// add header, and write the result into response body
	resbody := fmt.Sprintf(`{"result": %v}`, resStr)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resbody))

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
	// if the request has the wrong method
	// or it does not have application/json header given it needs
	// throw error
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
