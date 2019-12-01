package handlers

import (
	R "INFO441-RoomReservation/servers/reservation/models/reservations"

	"github.com/streadway/amqp"
)

type HandlerContext struct {
	ReservationStore *R.ReservationStore
	RabbitConnection *amqp.Channel
	RabbitQueueName  string
}
