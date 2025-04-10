package dtos

type CardWithCountDTO struct {
	CardID   int    `json:"cardID"`
	CardName string `json:"cardName"`
	Count    int    `json:"count"`
}
