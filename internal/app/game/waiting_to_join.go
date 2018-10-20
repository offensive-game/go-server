package game

import (
	"fmt"
	"go-server/internal/app/config"
	"go-server/internal/app/models"
	"time"
)

func (m Manager) WaitingToJoin() {
	for true {
		select {
		case joined := <-m.Input:
			{
				if joined.Order() == config.ORDER_JOIN {
					full := m.newPlayerJoined(joined)
					if full {
						break
					}
				}
			}
		case <-time.After(1000 * time.Second):
			{
				break
			}
		}
	}

}

func (m Manager) newPlayerJoined(command Command) bool {
	joinCommand := command.(models.PlayerJoined)

	fmt.Println(joinCommand)
	return false
}

func (m Manager) timeoutForJoining() {
	fmt.Println("TIMEOUT")
}
