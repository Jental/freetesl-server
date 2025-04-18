package dtos

type DiscardPileStateDTO struct {
	Player   []*CardInstanceStateDTO `json:"player"`
	Opponent []*CardInstanceStateDTO `json:"opponent"`
}
