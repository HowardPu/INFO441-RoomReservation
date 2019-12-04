package main

import (
	H "INFO441-RoomReservation/servers/reservation/handlers"
	S "INFO441-RoomReservation/servers/reservation/store"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/streadway/amqp"
)

// pie is on table for supper
// the pie was on the table for supper

var server = "mssql441.c2mbdnajn2pb.us-east-1.rds.amazonaws.com"
var user = "admin"
var password = "info441ishard"
var database = "RoomReservation"
var port = "1433"

var signingKey = "JusticsFromAbove"

var connString = fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;port=%s", server, user, password, database, port)
var db, dbERR = sql.Open("mssql", connString)

var resStore = S.NewReservationStore(db)

var queueName string = "reservationQueue"

var ctx = H.HandlerContext{
	RabbitQueueName:  queueName,
	ReservationStore: resStore,
	//RabbitConnection: ch,
}

// main is the main entry point for the server
func main() {
	/* TODO: add code to do the following
	- Read the ADDR environment variable to get the address
	  the server should listen on. If empty, default to ":80"
	- Create a new mux for the web server.
	- Tell the mux to call your handlers.SummaryHandler function
	  when the "/v1/summary" URL path is requested.
	- Start a web server listening on the address you read from
	  the environment variable, using the mux you created as
	  the root handler. Use log.Fatal() to report any errors
	  that occur when trying to start the web server.
	*/

	var conn, err = amqp.Dial("amqp://guest:guest@rabbit:5672/")

	for err != nil {
		failOnError(err, "RabbitMQ conn failed: ")
		time.Sleep(1 * time.Second)
		conn, err = amqp.Dial("amqp://guest:guest@rabbit:5672/")
	}

	var ch, chanErr = conn.Channel()

	var _, queueErr = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if dbERR != nil {
		failOnError(dbERR, "SQL conn failed: ")
	}

	if err != nil {
		failOnError(err, "RabbitMQ conn failed: ")
	}

	if chanErr != nil {
		failOnError(chanErr, "Create Chan Failed: ")
	}

	if queueErr != nil {
		failOnError(queueErr, "Create Queue Failed: ")
	}

	ctx.RabbitConnection = ch

	addr := os.Getenv("ADDR")
	defer conn.Close()
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/room", ctx.RoomHandler)
	mux.HandleFunc("/v1/reserve", ctx.RoomReserveHandler)
	mux.HandleFunc("/v1/specificRoom", ctx.SpecificRoomHandler)
	mux.HandleFunc("/v1/equip", ctx.EquipmentHandler)
	mux.HandleFunc("/v1/issue", ctx.IssueHandler)
	mux.HandleFunc("/v1/roomUsedTime", ctx.GetUsedTimeHandler)

	log.Printf("server is listening at %s...", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
