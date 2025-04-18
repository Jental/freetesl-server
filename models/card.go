package models

import (
	dbEnums "github.com/jental/freetesl-server/db/enums"
)

type Card struct {
	ID       int
	Name     string
	Cost     int
	ClassID  byte
	Type     dbEnums.CardType
	Keywords []dbEnums.CardKeyword
}

type CreatureCard struct {
	Card
	Power  int
	Health int
	Races  []byte
}

type ActionCard struct {
	Card
}

var cards []Card = []Card{CreatureCard{}.Card}
