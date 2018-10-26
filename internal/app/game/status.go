package game

import (
	"database/sql"
	"go-server/internal/app/config"
	"go-server/internal/app/models"
	"go-server/internal/app/utils"
	"time"
)

func (m Manager) sendGameStatus() error {
	round, phase, deadline, err := selectCurrentRound(m.db, m.GameModel.Id)
	if err != nil {
		return err
	}
	players, err := m.selectPlayers()
	if err != nil {
		return err
	}

	gameStatus := models.GameStatus{
		GameId:        m.GameModel.Id,
		Phase:         phase,
		Round:         round,
		RoundDeadline: utils.ToMillisecondsTimestamp(deadline),
		Players:       players,
	}

	m.sendToAllExcept(models.WebsocketNotification{
		Type:    models.PHASE_ADVANCE_SUCCESS,
		Payload: gameStatus,
	}, config.ALL_PLAYERS)

	return nil
}

func selectCurrentRound(db *sql.DB, gameId int64) (int, string, time.Time, error) {
	stmt, err := db.Prepare("SELECT round, phase, deadline FROM rounds WHERE gameId = $1")
	if err != nil {
		return 0, "", time.Time{}, err
	}

	row := stmt.QueryRow(gameId)

	var round int
	var phase string
	var deadline time.Time

	err = row.Scan(&round, &phase, &deadline)
	if err != nil {
		return 0, "", time.Time{}, err
	}

	return round, phase, deadline, nil
}

func (m Manager) selectPlayers() ([]models.PlayerStatus, error) {
	status := make(map[int64]models.PlayerStatus)

	players := m.getPlayersSlice()
	for _, player := range players {
		status[player.PlayerId()] = models.PlayerStatus{
			Id:             player.PlayerId(),
			Color:          player.PlayerColor(),
			Name:           player.PlayerName(),
			Cards:          []string{},
			Lands:          []models.Land{},
			UnitsInReserve: config.INITIAL_NUMBER_OF_UNITS,
		}
	}

	stmt, err := m.db.Prepare(
		`SELECT b.land, b.playerId, b.troops 
		FROM board b INNER JOIN players p ON p.id = b.playerId
		WHERE p.gameId = $1`)

	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(m.GameModel.Id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var currentPlayer int64
		var currentLand string
		var currentTroops int

		err = rows.Scan(&currentLand, &currentPlayer, &currentTroops)
		if err != nil {
			return nil, err
		}

		playerStatus := status[currentPlayer]
		playerStatus.Lands = append(playerStatus.Lands, models.Land{
			Name:          currentLand,
			NumberOfUnits: currentTroops,
		})

	}

	result := make([]models.PlayerStatus, 0)
	for _, p := range status {
		result = append(result, p)
	}

	return result, nil
}
