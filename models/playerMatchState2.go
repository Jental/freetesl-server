package models

import "github.com/gorilla/websocket"

type PlayerMatchState2 struct {
	PlayerID       int
	Connection     *websocket.Conn
	Deck           []*CardInstance
	Hand           []*CardInstance
	DiscardPile    []*CardInstance
	Health         int
	Runes          uint8
	Mana           int
	MaxMana        int
	LeftLaneCards  []*CardInstance
	RightLaneCards []*CardInstance
}
