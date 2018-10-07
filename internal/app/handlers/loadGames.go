package handlers

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"go-server/internal/app/middleware"
	"go-server/internal/app/utils"
	"net/http"
)

type LoadGames struct {
	appContext middleware.AppContext
	tx *sql.Tx
	log *logrus.Entry
}

type GamesList struct {
	games []Game
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
		SELECT id, player_count, name, start_time FROM games
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
	for rows.Next() {
		game := Game{}
		err := rows.Scan(&game.GameId, &game.NumberOfPlayers, &game.Name, &game.StartTime)
		if err != nil && err != sql.ErrNoRows {
			lg.log.Error("Unable to retrieve response from db")
		}
		games.games = append(games.games, game)
	}

	utils.RespondOK(&res, games)

}
