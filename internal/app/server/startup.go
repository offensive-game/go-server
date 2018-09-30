package server

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"go-server/internal/app/config"
	"go-server/internal/app/handlers"
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

	setUpHandlers(db)
	log.Info(fmt.Sprintf("Starting server on port %s", cfg.Port))
	err := http.ListenAndServeTLS(cfg.Port, "offensive.local.crt", "offensive.local.key", nil)
	if err != nil {
		log.Fatal("Unable to run server", err)
	}
}

func Handler(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte("Up and running"))
}

func setUpHandlers(db *sql.DB) {
	// Signup Handler
	signup := handlers.Signup{ Db: db }
	http.Handle("/signup", signup)
}
