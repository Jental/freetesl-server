package main

import (
	"flag"
	"log"
	"net/http"

	appHandlers "github.com/jental/freetesl-server/app/handlers"
	"github.com/jental/freetesl-server/match"
	matchActions "github.com/jental/freetesl-server/match/actions"
	matchHandlers "github.com/jental/freetesl-server/match/handlers"
	matchInterceptors "github.com/jental/freetesl-server/match/interceptors"
	matchSenders "github.com/jental/freetesl-server/match/senders"
	"github.com/jental/freetesl-server/models"
	"github.com/jental/freetesl-server/models/enums"
	"github.com/jental/freetesl-server/services"
)

var addr = flag.String("addr", "localhost:8081", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	go services.StartPlayersActivityMonitoring()

	match.MatchMessageHandlerFn = matchHandlers.ProcessMatchMessage
	match.BackendEventHandlerFn = matchSenders.ProcessBackendEvent

	var guardInterceptor models.Interceptor = matchInterceptors.GuardInterceptor{}
	matchInterceptors.RegisterInterceptor(enums.InterceptorPointHitFaceBefore, &guardInterceptor)
	matchInterceptors.RegisterInterceptor(enums.InterceptorPointHitCardBefore, &guardInterceptor)
	var coverInterceptor models.Interceptor = matchInterceptors.CoverInterceptor{}
	matchInterceptors.RegisterInterceptor(enums.InterceptorPointHitCardBefore, &coverInterceptor)

	matchInterceptors.RegisterAllSpecialCardsInterceptors()

	matchActions.RegisterAllActions()
	matchActions.RegisterActionsForCards()

	http.Handle("POST /login", http.HandlerFunc(appHandlers.Login))
	http.Handle("POST /logout", appHandlers.AuthCheckMiddleware(appHandlers.ActivityLoggerMiddleware(http.HandlerFunc(appHandlers.Logout))))
	http.Handle("GET /players", appHandlers.AuthCheckMiddleware(http.HandlerFunc(appHandlers.GetPlayers)))
	http.Handle("GET /currentPlayerInfo", appHandlers.AuthCheckMiddleware(http.HandlerFunc(appHandlers.GetCurrentPlayerInfo)))
	http.Handle("POST /lookingForOpponentStart", appHandlers.AuthCheckMiddleware(appHandlers.ActivityLoggerMiddleware(http.HandlerFunc(appHandlers.StartLookingForOpponent))))
	http.Handle("POST /lookingForOpponentStop", appHandlers.AuthCheckMiddleware(appHandlers.ActivityLoggerMiddleware(http.HandlerFunc(appHandlers.StopLookingForOpponent))))
	http.Handle("GET /lookingForOpponentStatus", appHandlers.AuthCheckMiddleware(appHandlers.ActivityLoggerMiddleware(http.HandlerFunc(appHandlers.GetLookingForOpponentStatus))))
	http.Handle("POST /matchCreate", appHandlers.AuthCheckMiddleware(appHandlers.ActivityLoggerMiddleware(http.HandlerFunc(matchHandlers.MatchCreate))))
	http.Handle("GET /decks", appHandlers.AuthCheckMiddleware(appHandlers.ActivityLoggerMiddleware(http.HandlerFunc(appHandlers.GetDecks))))
	http.Handle("/ws", appHandlers.AuthCheckMiddleware(http.HandlerFunc(matchHandlers.ConnectAndJoinMatch)))

	log.Fatal(http.ListenAndServe(*addr, appHandlers.RequestLoggerMiddleware(http.DefaultServeMux)))
}
