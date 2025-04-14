package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/jental/freetesl-server/models/enums"
	"github.com/jental/freetesl-server/services"
)

func DeleteDeck(w http.ResponseWriter, req *http.Request) {
	var playerID int = -1
	var ctx = req.Context()
	contextVal := ctx.Value(enums.ContextKeyUserID)
	if contextVal == nil {
		log.Println("player id is not found in a context")
	} else {
		var ok bool = false
		playerID, ok = contextVal.(int)
		if !ok {
			log.Println("player id from a context has invalid type")
			playerID = -1 // to have a common error handling
		}
	}
	if playerID < 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var urlParameters = req.URL.Query()
	var deckIDStrs, exists = urlParameters["deckID"]
	if !exists || len(deckIDStrs) == 0 {
		log.Printf("[%d]: DeleteDeck: no deck id passed", playerID)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	deckID, err := strconv.Atoi(deckIDStrs[0])
	if err != nil {
		log.Printf("[%d]: DeleteDeck: deck id is expected to be an int number", playerID)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	services.DeleteDeck(&ctx, playerID, deckID)
}
