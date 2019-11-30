package handlers

import (
	"time"

	U "INFO441-RoomReservation/servers/gateway/models/users"
)

//TODO: define a session state struct for this web server
//see the assignment description for the fields you should include
//remember that other packages can only see exported fields!
type SessionState struct {
	Time time.Time
	User *U.User
}
