package handlers

import (
	"github.com/gorilla/websocket"
	"go-server/internal/app/config"
	"go-server/internal/app/game"
	"go-server/internal/app/middleware"
	"go-server/internal/app/models"
	"go-server/internal/app/utils"
	"net/http"
)

type WebsocketOpen struct {
	appContext middleware.AppContext
}

func (w *WebsocketOpen) SetAppContext(appContext middleware.AppContext) {
	w.appContext = appContext
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (w *WebsocketOpen) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		panic("can't accept websocket connection")
	}

	query := req.URL.Query()
	token, ok := query["token"]
	if !ok {
		panic("missing auth token")
	}

	tx := utils.GetTransactionFromContext(req)

	statement, err := tx.Prepare(`
		SELECT p.id, p.color, u.username, p.gameId 
		FROM sessions s 
		INNER JOIN players p ON p.userId = s.userId 
		INNER JOIN users u ON u.id = p.userId 
		WHERE s.token = $1
	`)

	if err != nil {
		panic(err)
	}

	row := statement.QueryRow(token[0])

	var player models.Human
	var gameId int64

	err = row.Scan(&player.Id, &player.Color, &player.Name, &gameId)
	if err != nil {
		panic(err)
	}

	gameManager, found := game.GamesDictionary[gameId]
	if !found {
		panic("Can't find game with id")
	}

	command := models.PlayerJoined{Command: config.ORDER_JOIN, Player: player, Connection: conn}

	gameManager.Players[player.PlayerId()] = player
	gameManager.Input <- command
}
