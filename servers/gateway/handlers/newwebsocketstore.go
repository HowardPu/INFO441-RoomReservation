package handlers

import (
	"INFO441-RoomReservation/servers/gateway/sessions"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

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
		case client := <-h.register:
			token := client.token
			_, found := h.clients[token]
			if found {
				log.Printf("This token has been registered")
			}
			h.clients[token] = client
		case client := <-h.unregister:
			token := client.token
			if _, ok := h.clients[token]; ok {
				delete(h.clients, token)
				close(client.send)
			}
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

// serveWs handles websocket requests from the peer.
func (ctx *HandlerContext) ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// parse the user information
	// if not found, throw unauthorized
	userStore := UserLite{}
	_, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, &userStore)
	if err != nil || userStore.ID <= 0 {
		log.Println("User doesn't Sign in")
		http.Error(w, "User doesn't Sign in", http.StatusUnauthorized)
		return
	}

	// get the auth auery
	// if DNE, throw authorized
	authTokenQuery := r.URL.Query()["auth"]
	if len(authTokenQuery) == 0 {
		log.Println("No auth query")
		http.Error(w, "No Auth Query", http.StatusUnauthorized)
		return
	}

	// get auth token
	// upgrade the connection to websocket connection
	// if error occurs, throw bad request
	authToken := authTokenQuery[0]

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Conn Err: %v\n", err)
		return
	}

	log.Println("Called")
	client := &Client{
		hub:   hub,
		conn:  conn,
		send:  make(chan []byte, 256),
		token: authToken,
	}
	hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
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
