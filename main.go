package main

import (
	"flag"
	"log"
	"net/http"

	appHandlers "github.com/jental/freetesl-server/app/handlers"
	"github.com/jental/freetesl-server/match"
	matchHandlers "github.com/jental/freetesl-server/match/handlers"
	matchSenders "github.com/jental/freetesl-server/match/senders"
	"github.com/jental/freetesl-server/services"
)

var addr = flag.String("addr", "localhost:8081", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	go services.StartPlayersActivityMonitoring()

	match.MatchMessageHandlerFn = matchHandlers.ProcessMatchMessage
	match.BackendEventHandlerFn = matchSenders.ProcessBackendEvent

	http.Handle("POST /login", http.HandlerFunc(appHandlers.Login))
	http.Handle("POST /logout", appHandlers.AuthCheckMiddleware(appHandlers.ActivityLoggerMiddleware(http.HandlerFunc(appHandlers.Logout))))
	http.Handle("GET /players", appHandlers.AuthCheckMiddleware(http.HandlerFunc(appHandlers.GetPlayers)))
	http.Handle("GET /currentPlayerInfo", appHandlers.AuthCheckMiddleware(http.HandlerFunc(appHandlers.GetCurrentPlayerInfo)))
	http.Handle("POST /lookingForOpponentStart", appHandlers.AuthCheckMiddleware(appHandlers.ActivityLoggerMiddleware(http.HandlerFunc(appHandlers.StartLookingForOpponent))))
	http.Handle("POST /lookingForOpponentStop", appHandlers.AuthCheckMiddleware(appHandlers.ActivityLoggerMiddleware(http.HandlerFunc(appHandlers.StopLookingForOpponent))))
	http.Handle("GET /lookingForOpponentStatus", appHandlers.AuthCheckMiddleware(appHandlers.ActivityLoggerMiddleware(http.HandlerFunc(appHandlers.GetLookingForOpponentStatus))))
	http.Handle("POST /matchCreate", appHandlers.AuthCheckMiddleware(appHandlers.ActivityLoggerMiddleware(http.HandlerFunc(matchHandlers.MatchCreate))))
	http.Handle("/ws", appHandlers.AuthCheckMiddleware(http.HandlerFunc(matchHandlers.ConnectAndJoinMatch)))

	log.Fatal(http.ListenAndServe(*addr, appHandlers.RequestLoggerMiddleware(http.DefaultServeMux)))
}
