package dtos

type ApplyActionToCardDTO struct {
	CardInstanceID         string `json:"cardInstanceID"`
	OpponentCardInstanceID string `json:"opponentCardInstanceID"`
}
