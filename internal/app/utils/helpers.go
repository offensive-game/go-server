package utils

import (
	"database/sql"
	"net/http"
	"strings"
	"time"
)

func GetTransactionFromContext (req *http.Request) *sql.Tx {
	tx := req.Context().Value("tx").(*sql.Tx)
	return tx
}

func GetUserFromContext(req *http.Request) User {
	user := req.Context().Value("user").(User)
	return user
}

func ContainsMethods (methodsList string, method string) bool {
	methods := strings.Split(methodsList, ",")

	for _, current := range methods {
		if current == method {
			return true
		}
	}

	return false
}

func ToMillisecondsTimestamp (convertingTime time.Time) int64 {
	return convertingTime.UnixNano() / int64(time.Millisecond)
}