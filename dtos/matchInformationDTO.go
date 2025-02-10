package dtos

type MatchInformationDTO struct {
	Player   *PlayerInformationDTO `json:"player"`
	Opponent *PlayerInformationDTO `json:"opponent"`
}
