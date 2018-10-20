package game

import (
	"github.com/gorilla/websocket"
	"go-server/internal/app/models"
	"sync"
)

type Manager struct {
	JoinGameMutex *sync.Mutex
	GameModel     models.GameModel
	Input         chan Command
	Sockets       map[int64]*websocket.Conn
}

type Command interface {
	Order() string
}

var GamesDictionary = make(map[int64]Manager)

func NewGame(id int64) Manager {
	sockets := make(map[int64]*websocket.Conn)
	newManager := Manager{JoinGameMutex: &sync.Mutex{}, Input: make(chan Command), Sockets: sockets}
	GamesDictionary[id] = newManager

	go newManager.Run()

	return newManager
}

func (m Manager) Run() {
	m.WaitingToJoin()
}
