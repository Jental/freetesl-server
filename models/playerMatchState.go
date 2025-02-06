package models

type PlayerMatchState struct {
	Deck    []CardInstanceDTO
	Hand    []CardInstanceDTO
	Health  int
	Runes   uint8
	Mana    int
	MaxMana int
	OwnTurn bool
}
