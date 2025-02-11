package models

type Deck struct {
	ID       int
	Name     string
	PlayerID int
	Cards    map[int]CardWithCount
}
