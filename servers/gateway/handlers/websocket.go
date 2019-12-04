package handlers

import (
	"INFO441-RoomReservation/servers/gateway/sessions"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage = 1

	// BinaryMessage denotes a binary data message.
	BinaryMessage = 2

	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage = -1

	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage = 9

	// PongMessage denotes a pong control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage = 10
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

	// upgrade the connection
	// and pass it to the hub
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Conn Err: %v\n", err)
		return
	}

	client := &Client{
		hub:   hub,
		conn:  conn,
		send:  make(chan []byte, 256),
		token: authToken,
	}

	hub.register <- client

	endLoop := make(chan bool)
	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump(endLoop)

	// listenting the time when the connection is closed
	// if closed, end the handshake
	<-endLoop
}
