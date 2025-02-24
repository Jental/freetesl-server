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

	http.Handle("POST /login", appHandlers.RequestLoggerMiddleware(appHandlers.ActivityLoggerMiddleware(http.HandlerFunc(appHandlers.Login))))
	http.Handle("GET /players", appHandlers.RequestLoggerMiddleware(appHandlers.AuthCheckMiddleware(appHandlers.ActivityLoggerMiddleware(http.HandlerFunc(appHandlers.GetPlayers)))))
	http.Handle("POST /startLookingForOpponent", appHandlers.RequestLoggerMiddleware(appHandlers.AuthCheckMiddleware(appHandlers.ActivityLoggerMiddleware(http.HandlerFunc(appHandlers.StartLookingForOpponent)))))
	http.Handle("POST /stopLookingForOpponent", appHandlers.RequestLoggerMiddleware(appHandlers.AuthCheckMiddleware(appHandlers.ActivityLoggerMiddleware(http.HandlerFunc(appHandlers.StopLookingForOpponent)))))
	http.Handle("/ws", appHandlers.AuthCheckMiddleware(http.HandlerFunc(matchHandlers.ConnectAndJoinMatch)))

	log.Fatal(http.ListenAndServe(*addr, nil))
}
