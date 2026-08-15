package server

import (
	"backend/db"
	"backend/models"
	"backend/web"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type pageData struct{ Name, Email, Error string }

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()
	mw := &MiddleWares{DB: s.DB}
	stack := CreateStack(mw.Logging, mw.Recovery)
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) { s.render(w, r, "home.html", pageData{}, http.StatusOK) })
	mux.HandleFunc("GET /login", s.loginPage)
	mux.HandleFunc("POST /login", s.login)
	mux.HandleFunc("GET /register", s.registerPage)
	mux.HandleFunc("POST /register", s.register)
	mux.Handle("GET /me", mw.Auth(http.HandlerFunc(s.mePage)))
	mux.Handle("POST /logout", mw.Auth(http.HandlerFunc(s.logout)))
	return stack(mux)
}

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

func (s *Server) render(w http.ResponseWriter, r *http.Request, page string, data pageData, status int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		if err := web.RenderContent(w, page, data); err != nil {
			slog.Error("rendering page fragment failed", "page", page, "err", err)
		}
		return
	}
	w.WriteHeader(status)
	if err := web.RenderBase(w, page, data); err != nil {
		slog.Error("rendering page failed", "page", page, "err", err)
	}
}

func (s *Server) redirect(w http.ResponseWriter, r *http.Request, location string) {
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", location)
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Redirect(w, r, location, http.StatusSeeOther)
}
