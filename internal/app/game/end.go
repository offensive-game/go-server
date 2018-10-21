package game

import (
	"go-server/internal/app/config"
	"go-server/internal/app/models"
)

func (m Manager) endGame() {
	m.logger.Info("Game ended - insufficient number of players")

	tx, err := m.db.Begin()
	if err != nil {
		m.logger.Error("Unable to start transaction to clear finished game")
		return
	}

	statement, err := tx.Prepare("DELETE FROM games WHERE id = $1")
	if err != nil {
		m.logger.Error("Unable to prepare statement to clear finished game")
		return
	}

	_, err = statement.Exec(m.GameModel.Id)
	if err != nil {
		m.logger.Error("Unable to delete finished game")
		return
	}

	notification := models.WebsocketNotification{
		Type: "GAME_ENDED_NO_PLAYERS_SUCCESS",
		Payload: nil,
	}
	m.sendToAllExcept(notification, config.ALL_PLAYERS)
}
