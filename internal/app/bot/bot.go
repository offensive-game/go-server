package bot

import (
	"github.com/sirupsen/logrus"
	"go-server/internal/app/models"
)

const COMMAND_KILL = "kill"

type Executor struct {
	Game   models.GameModel
	Bot    models.Bot
	Output <-chan models.Command
	logger *logrus.Entry
}

func BuildBot(game models.GameModel, bot models.Bot, output <-chan models.Command, logger *logrus.Entry) Executor {
	botExec := Executor{
		Game:   game,
		Bot:    bot,
		Output: output,
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

		if command == COMMAND_KILL {
			b.logger.Info("Bot is finishing up")
			return
		}
	}
}
