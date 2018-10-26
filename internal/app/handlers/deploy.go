package handlers

import (
	"database/sql"
	"encoding/json"
	"go-server/internal/app/game"
	"go-server/internal/app/middleware"
	"go-server/internal/app/models"
	"go-server/internal/app/utils"
	"net/http"
)

type DeployMessage struct {
	GameId   int64  `json:"game_id"`
	PlayerId int64  `json:"player_id"`
	Land     string `json:"land"`
}

type Deploy struct {
	appContext middleware.AppContext
	tx         *sql.Tx
}

func (d *Deploy) SetAppContext(appContext middleware.AppContext) {
	d.appContext = appContext
}

func (d *Deploy) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	d.tx = utils.GetTransactionFromContext(req)
	body := DeployMessage{}

	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		utils.RespondBadRequest(&res, "invalid params")
		return
	}

	gameManager, found := game.GamesDictionary[body.GameId]
	if !found {
		utils.RespondBadRequest(&res, "invalid game")
		return
	}

	human, err := d.getPlayerForId(body.PlayerId)
	if err != nil {
		utils.RespondBadRequest(&res, "invalid player id")
		return
	}

	command := models.Deploy{
		Land:    body.Land,
		Player:  human,
		Success: make(chan bool),
	}

	gameManager.Input <- command
	success := <- command.Success

	if success {
		utils.RespondOK(&res, body)
	} else {
		utils.RespondBadRequest(&res, "can't deploy unit on that territory")
	}
}

func (d *Deploy) getPlayerForId(playerId int64) (*models.Human, error) {
	var human models.Human

	stmt, err := d.tx.Prepare(`
		SELECT p.id, u.username, p.color, p.units_in_reserve
		FROM players p 
		INNER JOIN users u ON u.id = p.userId
		WHERE p.id = $1
	`)
	if err != nil {
		return &human, err
	}

	row := stmt.QueryRow(playerId)

	err = row.Scan(&human.Id, &human.Name, &human.Color, &human.UnitsInReserve)
	if err != nil {
		return &human, err
	}

	return &human, nil
}
