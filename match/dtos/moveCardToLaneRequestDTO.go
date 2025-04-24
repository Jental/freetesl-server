package dtos

type MoveCardToLaneRequestDTO struct {
	CardInstanceID          string  `json:"cardInstanceID"`
	LaneID                  byte    `json:"laneID"`
	CardInstanceToReplaceID *string `json:"cardInstanceToReplaceID"`
}
