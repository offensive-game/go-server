package bot

import (
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"go-server/internal/app/config"
	"go-server/internal/app/models"
	"math/rand"
)

type Executor struct {
	Game   models.GameModel
	Bot    models.Bot
	Output chan<- models.Command
	logger *logrus.Entry
	db     *sql.DB
}

func BuildBot(game models.GameModel, bot models.Bot, output chan<- models.Command, logger *logrus.Entry, db *sql.DB) Executor {
	botExec := Executor{
		Game:   game,
		Bot:    bot,
		Output: output,
		db:     db,
		logger: logger.WithFields(logrus.Fields{
			"playerId":   bot.PlayerId(),
			"playerName": bot.PlayerName(),
		}),
	}

	return botExec
}

func (b *Executor) Run() {
	b.logger.Info("Bot is starting")
	for true {
		command := <-b.Bot.Input
		order := command.Type

		if order == models.COMMAND_KILL {
			b.logger.Info("Bot is finishing up")
			return
		} else if order == models.PHASE_ADVANCE_SUCCESS {
			err := b.advancePhase(command)
			if err != nil {
				b.logger.Error(err)
			}
		}

	}
}

func (b *Executor) advancePhase(message models.WebsocketNotification) error {
	gameStatus := message.Payload.(models.GameStatus)
	phase := gameStatus.Phase
	var err error
	switch phase {
	case config.DEPLOY:
		{
			err = b.deploy()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *Executor) deploy() error {
	stmt, err := b.db.Prepare(`
		SELECT land FROM board WHERE playerId = $1
	`)
	if err != nil {
		return err
	}

	rows, err := stmt.Query(b.Bot.Id)
	if err != nil {
		return err
	}

	myLands := make([]string, 0)
	for rows.Next() {
		var currentLand string
		err = rows.Scan(&currentLand)
		if err != nil {
			return err
		}
		myLands = append(myLands, currentLand)
	}

	for i := 0; i < b.Bot.PlayerUnitsInReserve(); i++ {
		randomInt := rand.Intn(len(myLands))
		randomLand := myLands[randomInt]
		b.logger.Info(fmt.Sprintf("DEPLOYING ON %s random number %d", randomLand, randomInt))

		message := models.Deploy{
			Player: &b.Bot,
			Land:   randomLand,
		}

		b.Output <- message
	}

	b.Bot.SetPlayerUnitsInReserve(0)

	return nil
}
