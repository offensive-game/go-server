package handlers

import (
	"database/sql"
	"net/http"
)

type Me struct {
	Db *sql.DB
}

func (me *Me) ServeHTTP(res http.ResponseWriter, req *http.Request) {

}
