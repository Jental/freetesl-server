package dtos

type DrawCardToLaneRequestDTO struct {
	LaneID                  byte    `json:"laneID"`
	CardInstanceToReplaceID *string `json:"cardInstanceToReplaceID"`
}
