package main

import (
	"backend/models"
	"backend/server"
	"backend/utils"
)

func main() {
	utils.SetupLogger("logs")
	models.GetEnv() // just load the env vars and validate them
	server.StartServer()
}
