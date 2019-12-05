package models

// this file represents the necessary struct for
// reservation store to perform correctly

type Reservation struct {
	ID          int    `json:"id"`
	TranDate    string `json:"tranDate"`
	ReserveDate string `json:"reserveDate"`
	RoomName    string `json:"roomName"`
	BeginTime   int    `json:"beginTime"`
	EndTime     int    `json:"endTime"`
	RoomType    string `json:"roomType"`
}

type Room struct {
	ID       int    `json:"id"`
	RoomName string `json:"roomName"`
	Capacity int    `json:"capacity"`
	Floor    int    `json:"floor"`
	RoomType string `json:"roomType"`
}

type Issue struct {
	ID          int    `json:"id"`
	RoomName    string `json:"roomName"`
	CreateDate  string `json:"createDate"`
	ConfirmDate string `json:"confirmDate"`
	SolveDate   string `json:"solveDate"`
	Body        string `json:"body"`
}

type Equipment struct {
	RoomEquipID int    `json:"roomEquipID"`
	Name        string `json:"Name"`
}

type User struct {
	ID   int64  `json:"userID"`
	Name string `json:"userName"`
	Type string `json:"userType"`
}
