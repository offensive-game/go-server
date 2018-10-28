package game

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"go-server/internal/app/models"
	"sync"
)

type Manager struct {
	JoinGameMutex *sync.Mutex
	GameModel     models.GameModel
	Input         chan models.Command
	Players       map[int64]models.Player
	joined        int8
	logger        *log.Entry
	db            *sql.DB
}

var GamesDictionary = make(map[int64]Manager)

func NewGame(currentGame models.GameModel, db *sql.DB) Manager {
	players := make(map[int64]models.Player)

	newManager := Manager{
		GameModel:     currentGame,
		JoinGameMutex: &sync.Mutex{},
		Input:         make(chan models.Command, 2),
		db:            db,
		Players:       players,
	}
	newManager.logger = log.WithFields(log.Fields{"gameId": currentGame.Id})

	GamesDictionary[currentGame.Id] = newManager

	go newManager.Run()

	return newManager
}

func (m *Manager) Run() {
	var err error

	defer func() {
		for _, player := range m.Players {
			if player.PlayerType() == models.BotType {
				player.SendMessage(models.WebsocketNotification{Type: models.COMMAND_KILL})
			}
		}
		m.endGame()
	}()

	// Waiting to join
	m.WaitingToJoin()

	err = m.initializeMap()
	if err != nil {
		return
	}

	m.sendGameStartMessage()

	// Deployment
	m.Deployment()

}

func (m *Manager) sendToAllExcept(message models.WebsocketNotification, playerId int64) {
	m.logger.Debug("sendToAllExcept")
	for id, player := range m.Players {
		if id == playerId {
			continue
		}
		m.logger.WithField("message", message).Debug("sending message")
		player.SendMessage(message)
	}
}

func (m Manager) getPlayersSlice() []models.Player {
	playersSlice := make([]models.Player, 0)

	for _, player := range m.Players {
		playersSlice = append(playersSlice, player)
	}

	return playersSlice
}
