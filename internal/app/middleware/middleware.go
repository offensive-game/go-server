package middleware

import (
	"context"
	"database/sql"
	"github.com/sirupsen/logrus"
	"go-server/internal/app/utils"
	"net/http"
)

type AppContext struct {
	DB     *sql.DB
	Logger *logrus.Logger
}

type AppHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	SetAppContext(AppContext)
}

type HandlerDecorator func(http.Handler, AppContext) http.HandlerFunc

func (appContext AppContext) Chain(method string, handler AppHandler, decorators ...HandlerDecorator) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		numberOfDecorators := len(decorators)

		if !utils.ContainsMethods(method, req.Method) {
			return
		}

		handler.SetAppContext(appContext)

		tx, err := appContext.DB.Begin()
		if err != nil {
			utils.RespondServerError(&res, "Server error")
		}

		newContext := context.WithValue(req.Context(), "tx", tx)
		defer func() {
			if r := recover(); r != nil {
				err := tx.Rollback()
				if err != nil {
					panic(err)
				}
			} else {
				err := tx.Commit()
				if err != nil {
					panic(err)
				}
			}
		}()

		if numberOfDecorators == 0 {
			handler.ServeHTTP(res, req.WithContext(newContext))
		} else {
			current := http.Handler(handler)
			for i := len(decorators) - 1; i >= 0; i-- {
				decorator := decorators[i]
				current = decorator(current, appContext)
			}

			current.ServeHTTP(res, req.WithContext(newContext))
		}
	})
}

func WithUser(handler http.Handler, appContext AppContext) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("offensive-login")
		if err != nil {
			utils.RespondNotAuthorized(&res)
			return
		}
		token := cookie.Value

		statement, err := appContext.DB.Prepare("SELECT u.id, u.username FROM sessions s INNER JOIN users u ON u.id = s.userId WHERE s.token = $1")
		if err != nil {
			panic(err)
		}

		row := statement.QueryRow(token)
		var user utils.User

		err = row.Scan(&user.Id, &user.Username)
		if err != nil {
			utils.RespondNotAuthorized(&res)
			return
		}

		ctx := context.WithValue(req.Context(), "user", user)
		handler.ServeHTTP(res, req.WithContext(ctx))
	})
}
