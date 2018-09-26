package main

import (
	"fmt"
	"go-server/internal/app/config"
	"go-server/internal/app/server"
)

func main () {
	c := config.Config{ Port: ":8080" }
	server.StartUpServer(c)
	fmt.Println("Djordje Vukovic")
}