package dtos

type MoveCardToLaneRequestDTO struct {
	PlayerID       int    `json:"playerID"`
	CardInstanceID string `json:"cardInstanceID"`
	LaneID         byte   `json:"laneID"`
}
