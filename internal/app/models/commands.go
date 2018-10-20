package models

import "github.com/gorilla/websocket"

type PlayerJoined struct {
	Command    string
	Player     Player
	Connection *websocket.Conn
}

func (p PlayerJoined) Order() string {
	return p.Command
}
