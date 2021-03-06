package server

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"go-server/internal/app/config"
	"go-server/internal/app/handlers"
	"go-server/internal/app/middleware"
	"net/http"
	"os"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func StartUpServer(cfg config.Config) {
	dbConnectionString := fmt.Sprintf("host=%s port=5432 user=offensive password=testee dbname=offensive sslmode=disable", cfg.DbHost)
	db, dbError := sql.Open("postgres", dbConnectionString)
	if dbError != nil {
		log.Fatal("Unable to connect to Database")
	}
	defer db.Close()

	clearUp(db)
	setUpHandlers(db)

	log.Info(fmt.Sprintf("Starting server on port %s", cfg.Port))
	err := http.ListenAndServeTLS(cfg.Port, "offensive.local.crt", "offensive.local.key", nil)
	if err != nil {
		log.Fatal("Unable to run server", err)
	}
}

func setUpHandlers(db *sql.DB) {
	appContext := middleware.AppContext{DB: db, Logger: log.New()}

	// Health check
	http.Handle("/hc", appContext.Chain("GET", &handlers.HealthCheck{}))

	// WebSocket handler
	http.Handle("/ws", appContext.Chain("GET", &handlers.WebsocketOpen{}))

	// Account handlers
	http.Handle("/signup", appContext.Chain("POST", &handlers.Signup{}))
	http.Handle("/login", appContext.Chain("POST", &handlers.Login{}))
	http.Handle("/signout", appContext.Chain("POST", &handlers.Signout{}, middleware.WithUser))

	// User/games management handlers
	http.Handle("/me", appContext.Chain("GET", &handlers.Me{}, middleware.WithUser))
	http.Handle("/game", appContext.Chain("POST,GET", &handlers.GameRequests{DB: db}, middleware.WithUser))
	http.Handle("/game/", appContext.Chain("POST,GET", &handlers.GameRequests{DB: db}, middleware.WithUser))

	// Game endpoints
	http.Handle("/deploy", appContext.Chain("POST", &handlers.Deploy{}, middleware.WithUser))
}

func clearUp(db *sql.DB) {
	db.Exec("DELETE FROM games")
}
