package models

// Room represents the room information
type Room struct {
	ID         int64  `json:"roomID"`
	Name       string `json:"roomName"`
	Floor      int64  `json:"floor"`
	Capacity   int64  `json:"capacity"`
	RoomType   string `json:"roomType"`
	StatusType string `json:"roomStatus"`
}
