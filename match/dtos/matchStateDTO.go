package dtos

type MatchStateDTO struct {
	Player                      PlayerMatchStateDTO `json:"player"`
	Opponent                    PlayerMatchStateDTO `json:"opponent"`
	OwnTurn                     bool                `json:"ownTurn"`
	WaitingForOtherPlayerAction bool                `json:"waitingForOtherPlayerAction"`
}
