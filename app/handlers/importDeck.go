package handlers

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jental/freetesl-server/db/queries"
	"github.com/jental/freetesl-server/dtos"
	"github.com/jental/freetesl-server/models/enums"
	"github.com/jental/freetesl-server/services"
)

func ImportDeck(w http.ResponseWriter, req *http.Request) {
	var playerID int = -1
	ctx := req.Context()
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

	var decoder = json.NewDecoder(req.Body)
	var requestDTO dtos.ImportDeckDTO
	err := decoder.Decode(&requestDTO)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bytes, err := base64.StdEncoding.DecodeString(requestDTO.FileBase64)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	deckAddRequest, err := services.ImportDeckFromExcel(playerID, bytes)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = queries.AddDeck(&ctx, deckAddRequest)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
