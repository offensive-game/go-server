package models

import (
	"go-server/internal/app/config"
)

type Command interface {
	Order() string
}

type PlayerJoined struct {
	Player     Player
}

func (p PlayerJoined) Order() string {
	return config.ORDER_JOIN
}

type Deploy struct {
	Player  Player
	Land    string
	Success chan bool
}

func (d Deploy) Order() string {
	return config.ORDER_DEPLOY
}
