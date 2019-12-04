package handlers

import (
	U "INFO441-RoomReservation/servers/gateway/models/users"
	S "INFO441-RoomReservation/servers/gateway/sessions"
)

//TODO: define a handler context struct that
//will be a receiver on any of your HTTP
//handler functions that need access to
//globals, such as the key used for signing
//and verifying SessionIDs, the session store
//and the user store

type HandlerContext struct {
	SessionKey   string
	SessionStore *S.RedisStore
	UserStore    *U.MsSqlStore
}
