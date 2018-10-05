package middleware

import (
	"context"
	"database/sql"
	"github.com/sirupsen/logrus"
	"go-server/internal/app/models"
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

		if req.Method != method {
			return
		}

		handler.SetAppContext(appContext)

		if numberOfDecorators == 0 {
			handler.ServeHTTP(res, req)
		} else {
			current := http.Handler(handler)
			for i := len(decorators) - 1; i >= 0; i-- {
				decorator := decorators[i]
				current = decorator(current, appContext)
			}

			current.ServeHTTP(res, req)
		}
	})
}

func WithUser(handler http.Handler, appContext AppContext) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("test-cookie")
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
		var user models.User

		err = row.Scan(&user.Id, &user.Username)
		if err != nil {
			utils.RespondNotAuthorized(&res)
			return
		}

		ctx := context.WithValue(req.Context(), "user", user)
		handler.ServeHTTP(res, req.WithContext(ctx))
	})
}
