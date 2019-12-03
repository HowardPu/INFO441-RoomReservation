package main

import (
	H "INFO441-RoomReservation/servers/gateway/handlers"
	U "INFO441-RoomReservation/servers/gateway/models/users"
	S "INFO441-RoomReservation/servers/gateway/sessions"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis"

	_ "github.com/denisenkom/go-mssqldb"
)

// make sql connection
var server = "mssql441.c2mbdnajn2pb.us-east-1.rds.amazonaws.com"
var user = "admin"
var password = "info441ishard"
var database = "RoomReservation"
var port = "1433"

// signing key for session
var signingKey = "JusticsFromAbove"

var connString = fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;port=%s", server, user, password, database, port)
var db, dbERR = sql.Open("mssql", connString)

var msstore = U.NewMsSqlStore(db)

var redisaddr = "redisServer:6379"

var client = redis.NewClient(&redis.Options{
	Addr: redisaddr,
})

var redisStore = S.NewRedisStore(client, time.Hour)

var ctx = H.HandlerContext{
	SessionKey:   signingKey,
	UserStore:    msstore,
	SessionStore: redisStore,
	SocketStore:  H.NewSocketStore(),
}

func main() {
	log.Println(dbERR)
	addr := os.Getenv("ADDR")
	reserveAddr := os.Getenv("RESERVE")

	//reserveURLs := []string{reserveAddr}

	if len(addr) == 0 {
		addr = ":443"
	}

	//tlsKeyPath := os.Getenv("TLSKEY")
	//tlsCertPath := os.Getenv("TLSCERT")

	//log.Printf(tlsKeyPath)
	//log.Printf(tlsCertPath)

	reserveProxy := ctx.NewServiceProxy(reserveAddr)

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/users", ctx.UsersHandler)
	mux.HandleFunc("/v1/sessions", ctx.SessionsHandler)
	mux.HandleFunc("/v1/sessions/", ctx.SpecificSessionHandler)
	mux.HandleFunc("/v1/ws", ctx.WebsocketHandler)

	mux.Handle("/v1/room", reserveProxy)
	mux.Handle("/v1/reserve", reserveProxy)
	mux.Handle("/v1/specificRoom", reserveProxy)
	mux.Handle("/v1/equip", reserveProxy)
	mux.Handle("/v1/issue", reserveProxy)

	go ctx.StartListeningRabbitMQ()

	wrappedMux := H.NewCors(mux)

	log.Printf("server is listening at %s...", addr)
	//log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, wrappedMux))
	log.Fatal(http.ListenAndServe(addr, wrappedMux))
}
