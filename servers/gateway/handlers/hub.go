package handlers

import (
	"log"

	"github.com/streadway/amqp"
)

// rabbit db/ query information
var rabbitAddr = "amqp://guest:guest@rabbit:5672/"
var queueName = "reservationQueue"
var durable = true

// Hub represents a center manage of all websocket connections
type Hub struct {
	// Registered clients.
	clients map[string]*Client

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

// create a new hub
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		// if the hub receives a new connection
		// add the connection to the hub
		case client := <-h.register:
			token := client.token
			h.clients[token] = client
		// if the hub receives close connection message
		// remove the connection in the hub
		case client := <-h.unregister:
			token := client.token
			if _, ok := h.clients[token]; ok {
				delete(h.clients, token)
			}
		// if the hub receives message from rabbitMQ
		// write it to all connection
		case message := <-h.broadcast:
			for token := range h.clients {
				client := h.clients[token]
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, token)
				}
			}
		}
	}
}

func (h *Hub) StartListeningRabbitMQ() {
	// make the rabbit connection
	conn, err := amqp.Dial(rabbitAddr)
	failOnError(err, "Fail to connect RabbitMQ")

	// make the channel
	channel, chanErr := conn.Channel()
	failOnError(chanErr, "Fail to open channel")

	_, qErr := channel.QueueDeclare(
		queueName, // namd
		durable,   // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // argument
	)

	// get the message from the rabbitMQ
	failOnError(qErr, "Fail to connect to Query")
	rabbitChan, errorChan := channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	failOnError(errorChan, "Fail to create consumption error")

	defer conn.Close()
	defer channel.Close()

	forever := make(chan bool)

	// listen the message from rabbitMQ
	// and write those message to the clients
	go func() {
		for d := range rabbitChan {
			h.broadcast <- []byte(d.Body)
			d.Ack(true)
		}
	}()

	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s \n", msg, err)
	}
}
