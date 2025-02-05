package dtos

type PlayerMatchStateDTO struct {
	Deck   []CardInstanceDTO `json:"deck"`
	Hand   []CardInstanceDTO `json:"hand"`
	Health int               `json:"health"`
	Runes  uint8             `json:"runes"`
}
