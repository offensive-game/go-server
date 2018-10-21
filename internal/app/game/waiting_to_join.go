package game

import (
	"fmt"
	"go-server/internal/app/config"
	"go-server/internal/app/models"
	"time"
)

func (m Manager) WaitingToJoin() {
	m.logger.Info("WAITING TO JOIN")
	for true {
		select {
		case joined := <-m.Input:
			{
				m.logger.Debug("JOINED RECEIVED")
				if joined.Order() == config.ORDER_JOIN {
					full := m.newPlayerJoined(joined)
					if full {
						break
					}
				}
			}
		case <-time.After(100000 * time.Second):
			{
				m.timeoutForJoining()
				break
			}
		}
	}

}

func (m *Manager) newPlayerJoined(command Command) bool {
	m.logger.Debug("newPlayerJoined")
	joinCommand := command.(models.PlayerJoined)

	opponentJoinedMessage := models.WebsocketNotification{
		Type: "OPPONENT_JOINED_SUCCESS",
		Payload: joinCommand.Player,
	}
	m.sendToAllExcept(opponentJoinedMessage, joinCommand.Player.Id)
	m.joined++

	return m.joined == m.GameModel.PlayersCount
}

func (m Manager) timeoutForJoining() {
	fmt.Println("TIMEOUT")
}
