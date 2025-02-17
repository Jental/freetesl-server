package dtos

type DeckStateDTO struct {
	Player   []*CardInstanceStateDTO `json:"player"`
	Opponent []*CardInstanceStateDTO `json:"opponent"`
}
