package handlers

import (
	"INFO441-RoomReservation/servers/gateway/sessions"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

// var rabbitAddr = "amqp://guest:guest@rabbit:5672/"
var rabbitAddr = "amqp://guest:guest@localhost:5672/"
var queueName = "reservationQueue"
var durable = true

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

const (
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage = 1

	// BinaryMessage denotes a binary data message.
	BinaryMessage = 2

	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage = 8

	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage = 9

	// PongMessage denotes a pong control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage = 10
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s \n", msg, err)
	}
}

//TODO: add a handler that upgrades clients to a WebSocket connection
//and adds that to a list of WebSockets to notify when events are
//read from the RabbitMQ server. Remember to synchronize changes
//to this list, as handlers are called concurrently from multiple
//goroutines.

func (ctx *HandlerContext) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	// parse the user information
	// if not found, throw unauthorized
	userStore := UserLite{}
	_, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, &userStore)
	if err != nil || userStore.ID <= 0 {
		http.Error(w, "User doesn't Sign in", http.StatusUnauthorized)
		return
	}

	// get the auth auery
	// if DNE, throw authorized
	authTokenQuery := r.URL.Query()["auth"]
	if len(authTokenQuery) == 0 {
		http.Error(w, "No Auth Query", http.StatusUnauthorized)
		return
	}

	// get auth token
	// upgrade the connection to websocket connection
	// if error occurs, throw bad request
	authToken := authTokenQuery[0]
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Fail to initialize websocket connection %v \n", err), http.StatusBadRequest)
		return
	}

	// create the websocket connection into websocket store
	// if error occurs, throw bad request
	createErr := ctx.SocketStore.AddNewConnection(authToken, conn)
	if createErr != nil {
		http.Error(w, fmt.Sprintf("Fail to create websocket connection %v \n", createErr), http.StatusBadRequest)
		return
	}
}

//TODO: start a goroutine that connects to the RabbitMQ server,
//reads events off the queue, and broadcasts them to all of
//the existing WebSocket connections that should hear about
//that event. If you get an error writing to the WebSocket,
//just close it and remove it from the list
//(client went away without closing from
//their end). Also make sure you start a read pump that
//reads incoming control messages, as described in the
//Gorilla WebSocket API documentation:
//http://godoc.org/github.com/gorilla/websocket

func (ctx *HandlerContext) StartListeningRabbitMQ() {
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

	// create thread to listen client message
	go ctx.ListeningClientMessage()

	// listen the message from rabbitMQ
	// and write those message to the clients
	go func() {
		for d := range rabbitChan {
			ctx.SocketStore.Lock.Lock()
			processErr := ctx.SocketStore.ProcessMessage(d.Body)
			if processErr != nil {
				log.Printf("Process Message Error: %v\n", processErr)
			}
			d.Ack(true)
			ctx.SocketStore.Lock.Unlock()
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	<-forever
}

func (ctx *HandlerContext) ListeningClientMessage() {
	// listen the spot when the websocket connection is created
	// and begin listen the client message
	for {
		auth := <-ctx.SocketStore.AuthChan
		go ctx.EndClientConnection(auth)
	}
}

func (ctx *HandlerContext) EndClientConnection(authToken string) {
	// get the websocket connection
	conn, found := ctx.SocketStore.Connections[authToken]
	if !found {
		log.Println("Connection Nor Found for this auth token")
		return
	}

	for {
		// listen the client message
		// if send close message
		// close the message
		messageType, _, err := conn.ReadMessage()
		if err != nil || messageType == CloseMessage {
			ctx.SocketStore.RemoveConnection(authToken)
			return
		}
	}
}
