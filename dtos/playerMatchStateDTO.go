package dtos

type PlayerMatchStateDTO struct {
	Deck []CardInstanceDTO `json:"deck"`
	Hand []CardInstanceDTO `json:"hand"`
}
