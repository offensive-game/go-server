package main

import (
	"../../internal/app/config"
	"../../internal/app/server"
)

func main() {
	conf := config.Config{ Port:"8080"}
	server.StartUpServer(conf)
}