package handlers

import (
	S "INFO441-RoomReservation/servers/reservation/store"

	"github.com/streadway/amqp"
)

type HandlerContext struct {
	ReservationStore *S.ReservationStore
	RabbitConnection *amqp.Channel
	RabbitQueueName  string
}
