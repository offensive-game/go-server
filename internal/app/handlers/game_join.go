package handlers

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"go-server/internal/app/config"
	"go-server/internal/app/game"
	"go-server/internal/app/middleware"
	"go-server/internal/app/models"
	"go-server/internal/app/utils"
	"net/http"
)

type JoinGame struct {
	appContext middleware.AppContext
	gameToJoin string
	tx         *sql.Tx
	log        *logrus.Entry
	user       utils.User
	DB         *sql.DB
}

func (g *JoinGame) SetAppContext(appContext middleware.AppContext) {
	g.appContext = appContext
}

func (g *JoinGame) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	g.appContext.Logger.Info("Handling JOIN GAME request")
	g.tx = utils.GetTransactionFromContext(req)
	g.user = utils.GetUserFromContext(req)

	g.log = g.appContext.Logger.WithFields(logrus.Fields{
		"userId": g.user.Id,
		"gameId": g.gameToJoin,
	})

	currentGame, err := g.getGame()
	if err != nil {
		utils.RespondBadRequest(&res, err.Error())
		return
	}

	gameManager, exists := game.GamesDictionary[currentGame.Id]
	if !exists {
		gameManager = game.NewGame(currentGame, g.DB)
	}

	gameManager.JoinGameMutex.Lock()
	defer gameManager.JoinGameMutex.Unlock()

	players, err := g.getPlayersForGame()
	if err != nil {
		utils.RespondBadRequest(&res, err.Error())
		return
	}

	if alreadyJoined(players, g.user.Username) {
		g.respondWithAlreadyJoined(&res, currentGame, players)
	} else {
		g.addNewPlayerToGame(&res, currentGame, players)
	}
}

func (g *JoinGame) respondWithAlreadyJoined(res *http.ResponseWriter, currentGame models.GameModel, players []models.Player) {
	var me models.Player
	for _, player := range players {
		if player.PlayerName() == g.user.Username {
			me = player
			break
		}
	}

	g.buildResponse(res, currentGame, me.PlayerColor(), me.PlayerId(), players)
}

func (g *JoinGame) addNewPlayerToGame(res *http.ResponseWriter, currentGame models.GameModel, players []models.Player) {
	if len(players) == config.MAX_NUMBER_PLAYERS-1 {
		utils.RespondBadRequest(res, "all players has already joined")
		return
	}

	assignedColors := make([]string, len(players))
	for i, player := range players {
		assignedColors[i] = player.PlayerColor()
	}
	newColor, err := utils.GetRandomColor(assignedColors[:])
	if err != nil {
		panic("Unable to find a new color")
	}

	newId, err := g.addNewPlayer(currentGame, newColor)

	players = append(players, models.Human{Name: g.user.Username, Id: newId, Color: newColor,})

	g.buildResponse(res, currentGame, newColor, newId, players)
}

func (g *JoinGame) getGame() (models.GameModel, error) {
	var currentGame models.GameModel

	statement, err := g.tx.Prepare(`
		SELECT id, players_count, name, start_time FROM games WHERE id = $1
	`)

	if err != nil {
		return currentGame, err
	}

	row := statement.QueryRow(g.gameToJoin)

	err = row.Scan(&currentGame.Id, &currentGame.PlayersCount, &currentGame.Name, &currentGame.StartTime)
	if err != nil {
		return currentGame, err
	}

	return currentGame, nil
}

func (g *JoinGame) getPlayersForGame() ([]models.Player, error) {
	players := make([]models.Player, 0, config.MAX_NUMBER_PLAYERS)

	statement, err := g.tx.Prepare(`
		SELECT p.id, p.color, u.username FROM players p INNER JOIN users u ON u.id = p.userId WHERE p.gameId = $1
	`)

	if err != nil {
		return players, err
	}

	rows, err := statement.Query(g.gameToJoin)
	if err != nil {
		return players, err
	}

	for rows.Next() {
		var player models.Human

		err = rows.Scan(&player.Id, &player.Color, &player.Name)

		if err != nil {
			return players, err
		}

		players = append(players, player)
	}

	return players, nil
}

func (g *JoinGame) addNewPlayer(game models.GameModel, color string) (int64, error) {
	statement, err := g.tx.Prepare(`
		INSERT INTO players (userId, gameId, color) VALUES ($1, $2, $3) RETURNING id
	`)

	if err != nil {
		return 0, err
	}

	row := statement.QueryRow(g.user.Id, game.Id, color)

	var newId int64
	err = row.Scan(&newId)

	if err != nil {
		return 0, err
	}

	return newId, nil
}

func (g *JoinGame) buildResponse(res *http.ResponseWriter, currentGame models.GameModel, myColor string, myId int64, players []models.Player) {
	response := models.JoinGameResponse{
		GameId:          currentGame.Id,
		StartTime:       utils.ToMillisecondsTimestamp(currentGame.StartTime),
		Name:            currentGame.Name,
		NumberOfPlayers: currentGame.PlayersCount,
		Color:           myColor,
		PlayerId:        myId,
		Players:         players,
	}

	utils.RespondOK(res, response)
}

func alreadyJoined(players []models.Player, username string) bool {
	for _, player := range players {
		if player.PlayerName() == username {
			return true
		}
	}

	return false
}
