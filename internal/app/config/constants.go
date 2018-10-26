package config

const MAX_NUMBER_PLAYERS = 6

func COLORS() [MAX_NUMBER_PLAYERS]string {
	return [...]string{"red", "green", "blue", "yellow", "purple", "brown"}
}

const (
	DEPLOY = "deployment"
	ATTACK = "attack"
	WARS   = "wars"
	MOVE   = "move"
)

const ALL_PLAYERS = -1
const INITIAL_NUMBER_OF_UNITS = 20

const DEPLOYMENT_DURATION = 60

const ORDER_JOIN = "ORDER_JOIN"
const ORDER_DEPLOY = "ORDER_DEPLOY"