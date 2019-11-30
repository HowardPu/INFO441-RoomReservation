package handlers

import (
	R "INFO441-RoomReservation/servers/reservation/models/reservations"
)

type HandlerContext struct {
	ReservationStore *R.MsSqlStore
}
