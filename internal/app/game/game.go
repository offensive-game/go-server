package game

import "sync"

type Manager struct {
	JoinGameMutex *sync.Mutex
}

var GamesDictionary = make(map[int64]Manager)

func NewGame(id int64) Manager {
	newManager := Manager{JoinGameMutex: &sync.Mutex{}}
	GamesDictionary[id] = newManager
	return newManager
}
