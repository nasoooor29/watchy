package server

import (
	"backend/db"
	"backend/models"
	"fmt"
	"log/slog"
	"net/http"
)

type Server struct {
	DB *db.SQL
}

func StartServer() {
	env := models.GetEnv()
	// create db
	db, err := db.GetDB()
	if err != nil {
		slog.Error("", "err", err)
		panic(err)
	}
	defer db.Close()
	server := &Server{
		DB: db,
	}

	// create mux
	// start the server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%v:%v", env.Host, env.Port),
		Handler: server.RegisterRoutes(),
	}
	slog.Info("server started", "host", env.Host, "port", env.Port)

	httpServer.ListenAndServe()
}
