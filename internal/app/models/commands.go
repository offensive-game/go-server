package models

import (
	"github.com/gorilla/websocket"
	"go-server/internal/app/config"
)

type Command interface {
	Order() string
}

type PlayerJoined struct {
	Player     Player
	Connection *websocket.Conn
}

func (p PlayerJoined) Order() string {
	return config.ORDER_JOIN
}

type Deploy struct {
	Player Player
	Land   string
}

func (d Deploy) Order() string {
	return config.ORDER_DEPLOY
}
