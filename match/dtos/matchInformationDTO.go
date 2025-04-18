package dtos

import "github.com/jental/freetesl-server/dtos"

type MatchInformationDTO struct {
	Player          *dtos.PlayerInformationDTO `json:"player"`
	Opponent        *dtos.PlayerInformationDTO `json:"opponent"`
	HasRing         bool                       `json:"hasRing"`
	OpponentHasRing bool                       `json:"opponentHasRing"`
}
