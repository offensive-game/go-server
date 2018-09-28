package server

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"go-server/internal/app/config"
	"net/http"
	"os"
	"strings"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func StartUpServer(cfg config.Config) {
	http.HandleFunc("/", Handler)
	log.Info(fmt.Sprintf("Starting server on port %s", cfg.Port))
	err := http.ListenAndServeTLS(cfg.Port, "offensive.local.crt", "offensive.local.key", nil)
	if err != nil {
		log.Fatal("Unable to run server", err)
	}
}

func Handler(response http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()

	responseSlice := make([]string, 0)

	for key, val := range query {
		responseSlice = append(responseSlice, fmt.Sprintf("Key: %s and valie %s", key, val))
	}

	answer := "Hello, queries you set are: Djordje Vukovic"
	response.Write([]byte(fmt.Sprintf("%s\n%s", answer, strings.Join(responseSlice, "\n"))))
}
