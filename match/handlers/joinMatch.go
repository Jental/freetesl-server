package handlers

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/jental/freetesl-server/match"
	"github.com/jental/freetesl-server/match/senders"
	"github.com/jental/freetesl-server/models/enums"
)

func JoinMatch(playerID int, connection *websocket.Conn) {
	matchState, playerState, _, err := match.GetCurrentMatchState(playerID)
	if err != nil {
		fmt.Println(err)
		return
	}

	playerState.Connection = connection

	go senders.StartListeningBackendEvents(playerState, matchState)

	senders.SendAllCardsToPlayer(playerState)

	playerState.Events <- enums.BackendEventCardInstancesChanged
	playerState.Events <- enums.BackendEventOpponentCardInstancesChanged
	playerState.Events <- enums.BackendEventHandChanged
	playerState.Events <- enums.BackendEventOpponentHandChanged
	playerState.Events <- enums.BackendEventDeckChanged
	playerState.Events <- enums.BackendEventOpponentDeckChanged
	playerState.Events <- enums.BackendEventDiscardPileChanged
	playerState.Events <- enums.BackendEventOpponentDiscardPileChanged
}
