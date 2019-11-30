package models

// Room represents the room information
type Room struct {
	ID         int    `json:"roomID"`
	Name       string `json:"roomName"`
	Floor      int    `json:"floor"`
	Capacity   int    `json:"capacity"`
	RoomType   string `json:"roomType"`
	StatusType string `json:"roomStatus"`
}
