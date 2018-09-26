package server

import (
	"fmt"
	"go-server/internal/app/config"
	"net/http"
	"strings"
)

func StartUpServer(cfg config.Config) {
	http.HandleFunc("/", Handler)
	http.ListenAndServe(cfg.Port, nil)
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