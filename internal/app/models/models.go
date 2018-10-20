package models

import "time"

type Player struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type GameModel struct {
	Id           int64
	PlayersCount int8
	Name         string
	StartTime    time.Time
}

