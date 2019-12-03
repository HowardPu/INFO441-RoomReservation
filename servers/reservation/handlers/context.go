package handlers

import (
	M "INFO441-RoomReservation/servers/reservation/models"

	"github.com/streadway/amqp"
)

type HandlerContext struct {
	ReservationStore *M.ReservationStore
	RabbitConnection *amqp.Channel
	RabbitQueueName  string
}
