package game

import (
	"go-server/internal/app/config"
	"go-server/internal/app/models"
	"time"
)

func (m Manager) WaitingToJoin() bool {
	m.logger.Info("WAITING TO JOIN")

	waitTime := m.GameModel.StartTime.Sub(time.Now())
	for true {
		select {
		case joined := <-m.Input:
			{
				m.logger.Debug("JOINED RECEIVED")
				if joined.Order() == config.ORDER_JOIN {
					full := m.newPlayerJoined(joined)
					if full {
						return true
					}
				}
			}
		case <-time.After(waitTime):
			{
				return false
			}
		}
	}

	return false
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
