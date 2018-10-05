package handlers

import (
	"database/sql"
	"go-server/internal/app/middleware"
	"net/http"
)

type LoadGames struct {
	appContext middleware.AppContext
	tx *sql.Tx
}

func (l LoadGames) SetAppContext(appContext middleware.AppContext) {
	l.appContext = appContext
}

func (lg LoadGames) ServeHTTP (res http.ResponseWriter, req *http.Request) {

}
