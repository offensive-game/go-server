package game

import (
	"database/sql"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"go-server/internal/app/models"
	"sync"
)

type Manager struct {
	JoinGameMutex *sync.Mutex
	GameModel     models.GameModel
	Input         chan Command
	Sockets       map[int64]*websocket.Conn
	joined        int8
	logger        *log.Entry
	db            *sql.DB
}

type Command interface {
	Order() string
}

var GamesDictionary = make(map[int64]Manager)

func NewGame(currentGame models.GameModel, db *sql.DB) Manager {
	sockets := make(map[int64]*websocket.Conn)
	newManager := Manager{
		GameModel:     currentGame,
		JoinGameMutex: &sync.Mutex{},
		Input:         make(chan Command),
		Sockets:       sockets,
		db:            db,
	}
	newManager.logger = log.WithFields(log.Fields{"gameId": currentGame.Id})

	GamesDictionary[currentGame.Id] = newManager

	go newManager.Run()

	return newManager
}

func (m *Manager) Run() {
	playersJoined := m.WaitingToJoin()

	if !playersJoined {
		m.endGame()
		return
	}
}

func (m *Manager) sendToAllExcept(message interface{}, playerId int64) {
	m.logger.Debug("sendToAllExcept")
	for id, socket := range m.Sockets {
		if id != playerId {
			err := socket.WriteJSON(message)
			if err != nil {
				m.logger.Info(err)
			}
		}
	}
}
