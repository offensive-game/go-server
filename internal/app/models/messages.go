package models

const (
	GAME_ENDED_NO_PLAYERS_SUCCESS = "GAME_ENDED_NO_PLAYERS_SUCCESS"
	COMMAND_KILL                  = "COMMAND_KILL"
	PHASE_ADVANCE_SUCCESS         = "PHASE_ADVANCE_SUCCESS"
	OPPONENT_JOINED_SUCCESS       = "OPPONENT_JOINED_SUCCESS"
	GAME_START_SUCCESS            = "GAME_START_SUCCESS"
)

type WebsocketNotification struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type JoinGameResponse struct {
	GameId          int64    `json:"game_id"`
	StartTime       int64    `json:"start_time"`
	Name            string   `json:"name"`
	NumberOfPlayers int8     `json:"number_of_players"`
	Color           string   `json:"color"`
	PlayerId        int64    `json:"player_id"`
	Players         []Player `json:"players"`
}

type PlayerStatus struct {
	Id             int64    `json:"id"`
	Name           string   `json:"name"`
	Color          string   `json:"color"`
	Lands          []Land   `json:"lands"`
	Cards          []string `json:"cards"`
	UnitsInReserve int      `json:"units_in_reserve"`
}

type GameStatus struct {
	GameId        int64          `json:"game_id"`
	Phase         string         `json:"phase"`
	Round         int            `json:"round"`
	RoundDeadline int64          `json:"round_deadline"`
	Players       []PlayerStatus `json:"players"`
}
