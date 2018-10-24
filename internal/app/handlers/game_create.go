package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"go-server/internal/app/middleware"
	"go-server/internal/app/utils"
	"net/http"
	"time"
)

type CreateGame struct {
	appContext middleware.AppContext
	log *logrus.Entry
}

type Game struct {
	Name string `json:"name"`
	NumberOfPlayers int8 `json:"number_of_players"`
	WaitTime int16 `json:"wait_time,omitempty"`
	StartTime int64 `json:"start_time,omitempty"`
	GameId int64 `json:"id"`
}

func (cg *CreateGame) SetAppContext (appContext middleware.AppContext) {
	cg.appContext = appContext
}

func (cg *CreateGame) ServeHTTP (res http.ResponseWriter, req *http.Request) {
	cg.appContext.Logger.Info("Handling CREATE GAME request")
	tx := utils.GetTransactionFromContext(req)
	body := Game{}

	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		cg.appContext.Logger.Error("Invalid params in CREATE GAME request")
		utils.RespondBadRequest(&res, "Invalid params for game creation")
		return
	}

	cg.log = cg.appContext.Logger.WithFields(logrus.Fields{
		"userId": utils.GetUserFromContext(req).Id,
		"request": body,
	})

	newGame := cg.createNewGame(tx, body)
	utils.RespondOK(&res, newGame)
}

func (cg CreateGame) createNewGame (tx *sql.Tx, newGame Game) Game {
	statement, err := tx.Prepare(`
		INSERT INTO games (name, players_count, start_time)
		VALUES ($1, $2, $3)
		RETURNING id; 
	`)

	if err != nil {
		cg.log.Error("Cart create INSERT INTO GAME Statement")
		panic(nil)
	}

	duration := fmt.Sprintf("%ds", newGame.WaitTime)
	startDuration, err := time.ParseDuration(duration)
	if err != nil {
		cg.log.Error(fmt.Sprintf("Cant create duration %s", duration))
		panic(err)
	}

	gameStart := time.Now().UTC().Add(startDuration)
	row := statement.QueryRow(newGame.Name, newGame.NumberOfPlayers, gameStart)

	var newId int64
	err = row.Scan(&newId)
	if err != nil {
		cg.log.Error("Cant read returning if from insertion in game table")
		panic(err)
	}

	newGame.WaitTime = 0
	newGame.StartTime = utils.ToMillisecondsTimestamp(gameStart)
	newGame.GameId = newId

	return newGame
}