package models

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/jental/freetesl-server/common"
)

type Match struct {
	Id          uuid.UUID
	Connection0 *websocket.Conn
	Connection1 *websocket.Conn
	Player0ID   common.Maybe[int]
	Player1ID   common.Maybe[int]
}
