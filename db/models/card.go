package models

import "github.com/jental/freetesl-server/db/enums"

type Card struct {
	ID       int
	Name     string
	Power    int
	Health   int
	Cost     int
	ClassID  byte
	Type     enums.CardType
	Keywords []enums.CardKeyword
	Races    []byte
	Effects  []CardEffect
}
