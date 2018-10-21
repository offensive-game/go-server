package config

const MAX_NUMBER_PLAYERS = 6

func COLORS() [MAX_NUMBER_PLAYERS]string {
	return [...]string {"red", "green", "blue", "yellow", "purple", "brown"}
}

const ALL_PLAYERS = -1

const ORDER_JOIN = "JOIN"