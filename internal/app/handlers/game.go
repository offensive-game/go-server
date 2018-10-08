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

type GameRequests struct {
	appContext middleware.AppContext
}

func (gr *GameRequests) SetAppContext (appContext middleware.AppContext) {
	gr.appContext = appContext
}

func (gr *GameRequests) ServeHTTP (res http.ResponseWriter, req *http.Request) {
	var handler middleware.AppHandler
	if req.Method == "POST" {
		handler = &CreateGame{}
	} else if req.Method == "GET" {
		handler = &LoadGames{}
	}

	if handler == nil {
		panic("No handler for specified request")
	}

	handler.SetAppContext(gr.appContext)
	handler.ServeHTTP(res, req)
}

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

	gameStart := time.Now().Add(startDuration)
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

type LoadGames struct {
	appContext middleware.AppContext
	tx *sql.Tx
	log *logrus.Entry
}

type GamesList struct {
	Games []Game `json:"games"`
}

func (lg *LoadGames) SetAppContext(appContext middleware.AppContext) {
	lg.appContext = appContext
}

func (lg *LoadGames) ServeHTTP (res http.ResponseWriter, req *http.Request) {
	lg.appContext.Logger.Info("Handling LOAD GAMES request")

	lg.log = lg.appContext.Logger.WithFields(logrus.Fields{
		"userId": utils.GetUserFromContext(req).Id,
	})

	tx := utils.GetTransactionFromContext(req)

	statement, err := tx.Prepare(`
		SELECT id, players_count, name, start_time FROM games
	`)

	if err != nil {
		lg.log.Error("Cant get games in progress")
		panic(err)
	}

	rows, err := statement.Query()
	defer func() {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}()

	if err != nil {
		lg.log.Error("Unable to execute query to fetch games in progress")
		panic(err)
	}

	games := GamesList{}
	games.Games = []Game{}
	for rows.Next() {
		game := Game{}

		var gameStartTime time.Time
		err := rows.Scan(&game.GameId, &game.NumberOfPlayers, &game.Name, &gameStartTime)
		if err != nil && err != sql.ErrNoRows {
			lg.log.Error("Unable to retrieve response from db")
			break
		}
		game.StartTime = utils.ToMillisecondsTimestamp(gameStartTime)
		games.Games = append(games.Games, game)
	}

	utils.RespondOK(&res, games)

}