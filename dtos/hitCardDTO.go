package dtos

type HitCardDTO struct {
	PlayerID               int    `json:"playerID"`
	CardInstanceID         string `json:"cardInstanceID"`
	OpponentCardInstanceID string `json:"opponentCardInstanceID"`
}
