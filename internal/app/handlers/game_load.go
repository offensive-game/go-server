package handlers

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"go-server/internal/app/middleware"
	"go-server/internal/app/utils"
	"net/http"
	"time"
)

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