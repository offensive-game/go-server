package config

const MAX_NUMBER_PLAYERS = 6

func COLORS() [MAX_NUMBER_PLAYERS]string {
	return [...]string {"red", "green", "blue", "yellow", "purple", "brown"}
}

const NO_PLAYER_ID = -1

const ORDER_JOIN = "JOIN"