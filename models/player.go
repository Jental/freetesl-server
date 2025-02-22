package models

import "github.com/jental/freetesl-server/models/enums"

type Player struct {
	ID          int
	DisplayName string
	AvatarName  string
	State       enums.PlayerState
}
