package game

import (
	"database/sql"
	"fmt"
	"go-server/internal/app/config"
	"go-server/internal/app/models"
	"time"
)

func (m *Manager) Deployment() {
	var err error
	m.logger.Info("deployment phase starts")
	err = m.sendGameStatus()
	if err != nil {
		m.logger.Error(err)
	}
	for {
		select {
		case message := <-m.Input:
			{
				m.logger.Debug("deploy message received")
				if message.Order() == config.ORDER_DEPLOY {
					m.deployUnit(message.(models.Deploy))
				}
			}
		case <-time.After(config.DEPLOYMENT_DURATION * time.Second):
			{
				m.logger.Info("deployment phase time-out")
				return
			}
		}
	}
}

func (m *Manager) deployUnit(deploy models.Deploy) {
	tx, err := m.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			m.logger.Fatal("Rolling back transaction in deployUnit")
			deploy.Success <- false
			return
		}
		deploy.Success <- true
		tx.Commit()
	}()

	if deploy.Player.PlayerUnitsInReserve() == 0 {
		m.logger.Info("can't deploy unit because there are no units in reserve")
		return
	}

	err = m.decreaseUnitsInReserve(tx, deploy.Player)
	if err != nil {
		return
	}

	err = m.increaseUnitsOnBoard(tx, deploy.Player, deploy.Land)
	if err != nil {
		return
	}

	for _, player := range m.Players {
		if player.PlayerId() == deploy.Player.PlayerId() {
			player.SetPlayerUnitsInReserve(player.PlayerUnitsInReserve() - 1)
			m.logger.Info(fmt.Sprintf("Units deployed, left %d", player.PlayerUnitsInReserve()))
		}
	}
}

func (m *Manager) decreaseUnitsInReserve(tx *sql.Tx, player models.Player) error {
	stmt, err := tx.Prepare(`UPDATE players SET units_in_reserve = units_in_reserve - 1 WHERE id = $1`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(player.PlayerId())
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) increaseUnitsOnBoard(tx *sql.Tx, player models.Player, land string) error {
	stmt, err := tx.Prepare(`UPDATE board SET troops = troops + 1 WHERE playerId = $1 AND land = $2`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(player.PlayerId(), land)
	if err != nil {
		return err
	}

	return nil
}
