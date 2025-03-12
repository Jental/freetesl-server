package handlers

import (
	"log"

	"github.com/google/uuid"
	"github.com/jental/freetesl-server/dtos"
	"github.com/jental/freetesl-server/models"
	"github.com/jental/freetesl-server/models/enums"
	"github.com/mitchellh/mapstructure"
)

func ProcessMatchMessage(playerID int, message models.PartiallyParsedMessage) error {
	log.Printf("[%d]: processing message: '%s'(%t)", playerID, message.Method, message.Method == "")

	switch message.Method {
	case "endTurn":
		go EndTurn(playerID)
	case "moveCardToLane":
		var dto dtos.MoveCardToLaneRequestDTO
		mapstructure.Decode(message.Body, &dto)
		cardInstanceID, err := uuid.Parse(dto.CardInstanceID)
		if err != nil {
			return err
		}
		laneID := enums.Lane(dto.LaneID)
		go MoveCardToLane(playerID, cardInstanceID, laneID)
	case "hitFace":
		var dto dtos.HitFaceDTO
		mapstructure.Decode(message.Body, &dto)
		cardInstanceID, err := uuid.Parse(dto.CardInstanceID)
		if err != nil {
			return err
		}
		go HitFace(playerID, cardInstanceID)
	case "hitCard":
		var dto dtos.HitCardDTO
		mapstructure.Decode(message.Body, &dto)
		cardInstanceID, err := uuid.Parse(dto.CardInstanceID)
		if err != nil {
			return err
		}
		opponentCardInstanceID, err := uuid.Parse(dto.OpponentCardInstanceID)
		if err != nil {
			return err
		}
		go HitCard(playerID, cardInstanceID, opponentCardInstanceID)
	case "concede":
		go Concede(playerID)
	}

	return nil
}
