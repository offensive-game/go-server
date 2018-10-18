package handlers

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"go-server/internal/app/middleware"
	"go-server/internal/app/utils"
	"net/http"
)

type JoinGame struct {
	appContext middleware.AppContext
	gameToJoin string
	tx *sql.Tx
	log *logrus.Entry
	user utils.User
}

func (g *JoinGame) SetAppContext(appContext middleware.AppContext) {
	g.appContext = appContext
}

func (g *JoinGame) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	g.appContext.Logger.Info("Handling JOIN GAME request")
	g.tx = utils.GetTransactionFromContext(req)
	g.user = utils.GetUserFromContext(req)

	g.log = g.appContext.Logger.WithFields(logrus.Fields{
		"userId": g.user.Id,
		"gameId": g.gameToJoin,
	})

	gameExists := g.gameExists()

	if !gameExists {

	}
}

func (g *JoinGame) gameExists() bool {
	statement, err := g.tx.Prepare(`
		SELECT COUNT(*) FROM games WHERE id = $1
	`)

	if err != nil {
		return false
	}

	row := statement.QueryRow(g.gameToJoin)

	var count int
	err = row.Scan(&count)
	if err != nil {
		return false
	}

	return count == 1
}