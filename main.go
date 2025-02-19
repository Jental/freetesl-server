package main

import (
	"flag"
	"log"
	"net/http"

	appHandlers "github.com/jental/freetesl-server/app/handlers"
	matchHandlers "github.com/jental/freetesl-server/match/handlers"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	http.HandleFunc("/login", appHandlers.Login)
	http.HandleFunc("/ws", matchHandlers.ConnectAndJoinMatch)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
