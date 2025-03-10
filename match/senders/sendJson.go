package senders

import (
	"fmt"

	"github.com/jental/freetesl-server/models"
)

func sendJson(playerState *models.PlayerMatchState, json map[string]interface{}) error {
	if playerState.Connection == nil {
		return fmt.Errorf("[%d]: sendJson: connection already closed", playerState.PlayerID)
	}

	playerState.SendMtx.Lock()
	defer playerState.SendMtx.Unlock()

	return playerState.Connection.WriteJSON(json)
}
