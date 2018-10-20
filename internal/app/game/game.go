package game

import (
	"go-server/internal/app/models"
	"sync"
)

type Manager struct {
	JoinGameMutex *sync.Mutex
	GameModel     models.GameModel
}

var GamesDictionary = make(map[int64]Manager)

func NewGame(id int64) Manager {
	newManager := Manager{JoinGameMutex: &sync.Mutex{}}
	GamesDictionary[id] = newManager
	return newManager
}

func Run () {

}
