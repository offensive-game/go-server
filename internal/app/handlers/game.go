package handlers

import (
	"go-server/internal/app/handlers/game_handlers"
	"go-server/internal/app/middleware"
	"net/http"
	"strings"
)

type GameRequests struct {
	appContext middleware.AppContext
}

func (gr *GameRequests) SetAppContext(appContext middleware.AppContext) {
	gr.appContext = appContext
}

func (gr *GameRequests) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var handler middleware.AppHandler
	if req.Method == "POST" {
		handler = &CreateGame{}
	} else if req.Method == "GET" {
		path := strings.Split(req.URL.Path[1:], "/")

		if len(path) == 1 {
			handler = &LoadGames{}
		} else {
			handler = &JoinGame{gameToJoin: path[1]}
		}

	}

	if handler == nil {
		panic("No handler for specified request")
	}

	handler.SetAppContext(gr.appContext)
	handler.ServeHTTP(res, req)
}
