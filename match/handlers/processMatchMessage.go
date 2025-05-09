package handlers

import (
	"log"

	"github.com/google/uuid"
	"github.com/jental/freetesl-server/match/dtos"
	"github.com/jental/freetesl-server/match/models"
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
		var cardInstanceToReplaceID *uuid.UUID
		if dto.CardInstanceToReplaceID != nil && *dto.CardInstanceToReplaceID != "" {
			id, err := uuid.Parse(*dto.CardInstanceToReplaceID)
			if err != nil {
				return err
			}
			cardInstanceToReplaceID = &id
		}
		laneID := enums.LanePosition(dto.LaneID)
		go MoveCardToLane(playerID, cardInstanceID, laneID, cardInstanceToReplaceID)
	case "drawCardToLane":
		var dto dtos.DrawCardToLaneRequestDTO
		mapstructure.Decode(message.Body, &dto)
		var cardInstanceToReplaceID *uuid.UUID
		if dto.CardInstanceToReplaceID != nil && *dto.CardInstanceToReplaceID != "" {
			id, err := uuid.Parse(*dto.CardInstanceToReplaceID)
			if err != nil {
				return err
			}
			cardInstanceToReplaceID = &id
		}
		laneID := enums.LanePosition(dto.LaneID)
		go DrawCardToLane(playerID, laneID, cardInstanceToReplaceID)
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
	case "drawCard":
		go DrawCard(playerID)
	case "applyCardToCard":
		var dto dtos.ApplyCardToCardDTO
		mapstructure.Decode(message.Body, &dto)
		cardInstanceID, err := uuid.Parse(dto.CardInstanceID)
		if err != nil {
			return err
		}
		opponentCardInstanceID, err := uuid.Parse(dto.OpponentCardInstanceID)
		if err != nil {
			return err
		}
		go ApplyCardToCard(playerID, cardInstanceID, opponentCardInstanceID)
	case "concede":
		go Concede(playerID)
	case "waitedUserActionsCompleted":
		go WaitedUserActionsCompleted(playerID)
	case "useRing":
		go UseRing(playerID)
	}

	return nil
}
