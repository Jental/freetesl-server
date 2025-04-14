package models

type AddDeckDbRequest struct {
	Name       string
	AvatarName string
	PlayerID   int
	Cards      map[int]int
}
