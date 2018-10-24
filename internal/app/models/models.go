package models

import "time"

type Land struct {
	Name          string `json:"name"`
	NumberOfUnits int    `json:"number_of_units"`
}

type GameModel struct {
	Id           int64
	PlayersCount int8
	Name         string
	StartTime    time.Time
}
