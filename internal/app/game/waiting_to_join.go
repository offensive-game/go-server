package game

import (
	"database/sql"
	"go-server/internal/app/bot"
	"go-server/internal/app/config"
	"go-server/internal/app/models"
	"go-server/internal/app/utils"
	"time"
)

func (m *Manager) WaitingToJoin() {
	m.logger.Info("WAITING TO JOIN")

	waitTime := m.GameModel.StartTime.Sub(time.Now().UTC())
	for true {
		select {
		case joined := <-m.Input:
			{
				m.logger.Debug("JOINED RECEIVED")
				if joined.Order() == config.ORDER_JOIN {
					m.newPlayerJoined(joined)
					if m.joined == m.GameModel.PlayersCount {
						return
					}
				}
			}
		case <-time.After(waitTime):
			{
				botsNeeded := m.GameModel.PlayersCount - m.joined
				for i := 0; i < int(botsNeeded); i++ {
					botPlayer := m.newBotJoined()
					m.Players[botPlayer.PlayerId()] = &botPlayer
					botRoutine := bot.BuildBot(m.GameModel, botPlayer, m.Input, m.logger, m.db)
					go botRoutine.Run()
				}
				return
			}
		}
	}
}

func (m *Manager) newBotJoined() models.Bot {
	m.logger.Debug("New bot joining")

	colorsAssigned := make([]string, 0)
	for _, player := range m.Players {
		colorsAssigned = append(colorsAssigned, player.PlayerColor())
	}
	newColor, err := utils.GetRandomColor(colorsAssigned)
	if err != nil {
		m.logger.Error("Unable to find new color")
		panic(err)
	}

	tx, err := m.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}

		tx.Commit()
	}()

	if err != nil {
		m.logger.Error("cant start transaction")
		panic(err)
	}

	stmt, err := tx.Prepare(`
		INSERT INTO players (userId, gameId, color, bot, units_in_reserve) VALUES ($1, $2, $3, $4, $5) RETURNING id
	`)
	if err != nil {
		m.logger.Error("cant prepare statement")
		panic(err)
	}

	row := stmt.QueryRow(nil, m.GameModel.Id, newColor, true, config.INITIAL_NUMBER_OF_UNITS)

	var botId int64
	err = row.Scan(&botId)
	if err != nil {
		m.logger.Error("Can't get id for new bot")
		panic(err)
	}

	name := "bot " + newColor
	newBot := models.Bot{
		Id:             botId,
		Color:          newColor,
		Name:           name,
		Input:          make(chan models.WebsocketNotification, 5),
		UnitsInReserve: config.INITIAL_NUMBER_OF_UNITS,
	}

	opponentJoinedMessage := models.WebsocketNotification{
		Type:    models.OPPONENT_JOINED_SUCCESS,
		Payload: newBot,
	}

	m.sendToAllExcept(opponentJoinedMessage, newBot.Id)

	return newBot
}

func (m *Manager) newPlayerJoined(command models.Command) {
	m.logger.Debug("newPlayerJoined")
	joinCommand := command.(models.PlayerJoined)

	m.Players[joinCommand.Player.PlayerId()] = joinCommand.Player

	opponentJoinedMessage := models.WebsocketNotification{
		Type:    models.OPPONENT_JOINED_SUCCESS,
		Payload: joinCommand.Player,
	}
	m.sendToAllExcept(opponentJoinedMessage, joinCommand.Player.PlayerId())
	m.joined++
}

func (m Manager) sendGameStartMessage() {
	playersList := make([]models.Player, 0, config.MAX_NUMBER_PLAYERS)

	for _, p := range m.Players {
		playersList = append(playersList, p)
	}

	playersJoined := models.JoinGameResponse{
		GameId:          m.GameModel.Id,
		StartTime:       utils.ToMillisecondsTimestamp(m.GameModel.StartTime),
		Name:            m.GameModel.Name,
		NumberOfPlayers: m.GameModel.PlayersCount,
		Players:         playersList,
	}

	for _, player := range playersList {
		playersJoined.PlayerId = player.PlayerId()
		playersJoined.Color = player.PlayerColor()

		notification := models.WebsocketNotification{
			Type:    models.GAME_START_SUCCESS,
			Payload: playersJoined,
		}

		player.SendMessage(notification)
	}
}

func (m *Manager) initializeMap() error {
	tx, err := m.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	if err != nil {
		m.logger.Error("can't start transaction in initializeMap")
		return err
	}

	err = m.createInitialRound(tx)
	if err != nil {
		return err
	}

	err = m.createInitialBoard(tx)
	if err != nil {
		return err
	}

	return nil
}

func (m Manager) createInitialRound(tx *sql.Tx) error {
	statement, err := tx.Prepare(`
		INSERT INTO rounds (round, phase, deadline, gameId) VALUES ($1, $2, $3, $4)
	`)

	if err != nil {
		return err
	}

	deadline := time.Now().UTC().Add(config.DEPLOYMENT_DURATION * time.Second)
	_, err = statement.Exec(1, "deployment", deadline, m.GameModel.Id)
	if err != nil {
		return err
	}

	return nil
}

func (m Manager) createInitialBoard(tx *sql.Tx) error {
	lands, err := m.getAllLands(tx)
	if err != nil {
		return err
	}

	players := m.getPlayersSlice()

	i := 0
	for _, land := range lands {
		stmt, err := tx.Prepare(`
			INSERT INTO board (land, playerId, troops) VALUES ($1, $2, 1)
		`)

		if err != nil {
			m.logger.Error("unable to create land for initializing board" + land.Name)
			panic(err)
		}
		_, err = stmt.Exec(land.Name, players[i].PlayerId())
		if err != nil {
			m.logger.Error("Cant insert land for initializing board")
			panic(err)
		}
		i = (i + 1) % len(players)
	}

	return nil
}

func (m Manager) getAllLands(tx *sql.Tx) ([]models.Land, error) {
	lands := make([]models.Land, 0)

	rows, err := tx.Query(`
		SELECT name FROM lands; 
	`)

	if err != nil {
		return lands, err
	}

	for rows.Next() {
		var land models.Land
		err = rows.Scan(&land.Name)
		if err != nil {
			return lands, err
		}
		lands = append(lands, land)
	}

	return lands, nil
}
