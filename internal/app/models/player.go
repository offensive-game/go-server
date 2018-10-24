package models

import "github.com/gorilla/websocket"

type PlayerType string

const (
	HumanType PlayerType = "HUMAN"
	BotType   PlayerType = "BOT"
)

type Player interface {
	PlayerId() int64
	PlayerName() string
	PlayerColor() string
	PlayerType() PlayerType
	SendMessage(WebsocketNotification)
}

type Human struct {
	Id     int64           `json:"id"`
	Name   string          `json:"name"`
	Color  string          `json:"color"`
	Socket *websocket.Conn `json:"-"`
}

func (h Human) PlayerId() int64 {
	return h.Id
}

func (h Human) PlayerName() string {
	return h.Name
}

func (h Human) PlayerColor() string {
	return h.Color
}

func (h Human) PlayerType() PlayerType {
	return HumanType
}

func (h Human) SendMessage(message WebsocketNotification) {
	h.Socket.WriteJSON(message)
}

type Bot struct {
	Id    int64       `json:"id"`
	Name  string      `json:"name"`
	Color string      `json:"color"`
	Input chan string `json:"-"`
}

func (b Bot) PlayerId() int64 {
	return b.Id
}

func (b Bot) PlayerName() string {
	return b.Name
}

func (b Bot) PlayerColor() string {
	return b.Color
}

func (b Bot) PlayerType() PlayerType {
	return BotType
}

func (b Bot) SendMessage(message WebsocketNotification) {
	b.Input <- message.Type
}
