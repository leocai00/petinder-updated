package main

import (
	"fmt"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/final-project-petinder/servers/gateway/handlers"
	"github.com/final-project-petinder/servers/gateway/indexes"
	"github.com/final-project-petinder/servers/gateway/models/users"
	"github.com/final-project-petinder/servers/gateway/sessions"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/streadway/amqp"
)

//main is the main entry point for the server
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
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":443"
	}

	sessionkey := os.Getenv("SESSIONKEY")
	if len(sessionkey) == 0 {
		sessionkey = "default"
	}

	redisaddr := os.Getenv("REDISADDR")
	redisdb := redis.NewClient(&redis.Options{
		Addr:     redisaddr,
		Password: "",
		DB:       0,
	})
	store := sessions.NewRedisStore(redisdb, time.Hour)

	tlscert := os.Getenv("TLSCERT")
	if len(tlscert) == 0 {
		log.Fatal("TLSCERT not set.")
	}

	tlskey := os.Getenv("TLSKEY")
	if len(tlskey) == 0 {
		log.Fatal("TLSKEY not set.")
	}

	sqlpwd := os.Getenv("MYSQL_ROOT_PASSWORD")
	if len(sqlpwd) == 0 {
		log.Fatal("MYSQL_ROOT_PASSWORD not set.")
	}

	dbaddr := os.Getenv("DBADDR")
	if len(dbaddr) == 0 {
		dbaddr = "127.0.0.1:3306"
	}

	dsn := fmt.Sprintf("root:%s@tcp(%s)/db", sqlpwd, dbaddr)
	database, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	db := users.NewMySQLStore(database)

	handler := &handlers.MyHandler{
		Key:          sessionkey,
		SessionStore: store,
		UserStore:    db,
		Trie:         indexes.NewTrie(),
		SocketStore:  handlers.NewSocketStore(),
	}

	err = handler.UserStore.Load(handler.Trie)
	if err != nil {
		handler.Trie = indexes.NewTrie()
	}

	messageaddr := os.Getenv("MESSAGEADDR")
	rabbitaddr := os.Getenv("RABBITADDR")
	petaddr := os.Getenv("PETADDR")

	conn, err := amqp.Dial("amqp://" + rabbitaddr + "/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q1, err := ch.QueueDeclare(
		"MsgQueue", // name matches what we used in our nodejs services
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare a queue")
	q2, err := ch.QueueDeclare(
		"PetQueue", // name matches what we used in our nodejs services
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q1.Name, // queue
		"",      // consumer
		false,   // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	failOnError(err, "Failed to register a consumer")
	petMsgs, err := ch.Consume(
		q2.Name, // queue
		"",      // consumer
		false,   // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	failOnError(err, "Failed to register a consumer")

	// Invoke a goroutine for handling control messages from this connection
	go handler.SocketStore.ProcessMessages(msgs)
	go handler.SocketStore.ProcessMessages(petMsgs)

	mux := http.NewServeMux()

	mux.Handle("/v1/channels", handler.NewServiceProxy(messageaddr))
	mux.Handle("/v1/channels/", handler.NewServiceProxy(messageaddr))
	mux.Handle("/v1/messages/", handler.NewServiceProxy(messageaddr))

	mux.Handle("/v1/pet", handler.NewServiceProxy(petaddr))
	mux.Handle("/v1/pet/", handler.NewServiceProxy(petaddr))

	mux.HandleFunc("/v1/users", handler.UsersHandler)
	mux.HandleFunc("/v1/users/", handler.SpecificUserHandler)
	mux.HandleFunc("/v1/sessions", handler.SessionsHandler)
	mux.HandleFunc("/v1/sessions/", handler.SpecificSessionHandler)

	mux.HandleFunc("/v1/ws", handler.WebSocketConnectionHandler)

	wrappedMux := handlers.NewResponseHeader(mux)

	log.Printf("server is listening at %s", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlscert, tlskey, wrappedMux))
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
