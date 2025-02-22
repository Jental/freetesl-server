package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jental/freetesl-server/dtos"
	"github.com/jental/freetesl-server/mappers"
	"github.com/jental/freetesl-server/services"
)

func GetPlayers(w http.ResponseWriter, req *http.Request) {
	var players, err = services.GetPlayers()
	if err != nil {
		log.Println(err)
		return
	}

	var responseDTO dtos.ListDTO[*dtos.PlayerInformationDTO] = dtos.ListDTO[*dtos.PlayerInformationDTO]{
		Items: mappers.MapToPlayerInformationDTOs(players),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseDTO)
}
