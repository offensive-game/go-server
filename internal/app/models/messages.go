package models

type PlayerModel struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
	Color string `json:"color"`
}

type JoinGameResponse struct {
	GameId int64 `json:"game_id"`
	StartTime int64 `json:"start_time"`
	Name string `json:"name"`
	NumberOfPlayers int8 `json:"number_of_players"`
	Color string `json:"color"`
	PlayerId int64 `json:"player_id"`
	Players []PlayerModel `json:"players"`
}