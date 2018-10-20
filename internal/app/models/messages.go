package models

type JoinGameResponse struct {
	GameId int64         `json:"game_id"`
	StartTime int64      `json:"start_time"`
	Name string          `json:"name"`
	NumberOfPlayers int8 `json:"number_of_players"`
	Color string         `json:"color"`
	PlayerId int64       `json:"player_id"`
	Players []Player     `json:"players"`
}

type WebsocketNotification struct {
	Type string `json:"type"`
	Payload interface{} `json:"payload"`
}