package utils

import (
	"database/sql"
	"go-server/internal/app/models"
	"net/http"
)

func GetTransactionFromContext (req *http.Request) *sql.Tx {
	tx := req.Context().Value("tx").(*sql.Tx)
	return tx
}

func GetUserFromContext(req *http.Request) models.User {
	user := req.Context().Value("user").(models.User)
	return user
}