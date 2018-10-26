package game

import (
	"database/sql"
	"go-server/internal/app/config"
	"go-server/internal/app/models"
)

func (m Manager) endGame() {
	m.logger.Info("Game ended - insufficient number of players")

	tx, err := m.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			m.logger.Error(err)
			m.logger.Fatal("Unable to clear game after ended" + string(m.GameModel.Id))
			return
		}
		tx.Commit()
	}()

	if err != nil {
		m.logger.Error("Unable to start transaction to clear finished game")
		return
	}

	err = clearUp(tx, m.GameModel.Id)
	if err != nil {
		return
	}

	notification := models.WebsocketNotification{
		Type: models.GAME_ENDED_NO_PLAYERS_SUCCESS,
		Payload: nil,
	}
	m.sendToAllExcept(notification, config.ALL_PLAYERS)
}

func clearUp (tx *sql.Tx, gameId int64) error {
	statement, err := tx.Prepare(`DELETE FROM players WHERE gameId = $1`)
	if err != nil {
		return err
	}
	_, err = statement.Exec(gameId)
	if err != nil {
		return err
	}

	statement, err = tx.Prepare("DELETE FROM games WHERE id = $1")
	if err != nil {
		return err
	}

	_, err = statement.Exec(gameId)
	if err != nil {
		return err
	}

	return nil
}
