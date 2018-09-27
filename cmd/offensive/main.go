package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"go-server/internal/app/config"
	"go-server/internal/app/server"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Port env variable is not supplied")
	}
	c := config.Config{Port: fmt.Sprintf(":%s", port)}
	server.StartUpServer(c)
}
