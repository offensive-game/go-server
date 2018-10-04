package server

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"go-server/internal/app/config"
	"go-server/internal/app/handlers"
	"go-server/pkg/middleware"
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

func setUpHandlers(db *sql.DB) {
	// Health check
	http.Handle("/hc", middleware.Chain("GET", http.HandlerFunc(hc)))

	// Signup Handler
	signup := handlers.Signup{Db: db}
	http.Handle("/signup", middleware.Chain("POST", signup))

	// Login Handler
	login := handlers.Login{Db: db}
	http.Handle("/login", middleware.Chain("POST", login))
}

func hc(res http.ResponseWriter, _ *http.Request) {
	_, err := res.Write([]byte("Alive"))
	if err != nil {
		panic(err)
	}
}
