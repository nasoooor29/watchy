package server

import (
	"backend/db"
	"backend/models"
	"backend/source"
	"backend/web"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()
	mw := &MiddleWares{DB: s.DB}
	stack := CreateStack(mw.Logging, mw.Recovery)
	mux.Handle("GET /public/", mw.cacheStaticFiles(http.StripPrefix("/public/", http.FileServer(http.Dir("web/public")))))
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		// if the url path is not / then return 404
		if r.URL.Path != "/" {
			s.render(w, r, http.StatusNotFound, "404.html")
			return
		}

		page, err := web.BuildLibraryPage(s.DB)
		if err != nil {
			slog.Error("failed to build library page", "err", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		s.render(w, r, http.StatusOK, "home.html", page)
	})

	mux.HandleFunc("GET /api/poster/{id}", s.GetPoster)
	mux.HandleFunc("GET /shows/{id}", s.GetShowPage)
	mux.HandleFunc("GET /shows/{id}/season/{season}", s.GetShowSeason)

	mux.HandleFunc("GET /login", s.RenderPageHandler("login.html"))
	mux.HandleFunc("POST /login", s.login)

	mux.HandleFunc("GET /register", s.RenderPageHandler("register.html"))
	mux.HandleFunc("POST /register", s.register)

	mux.Handle("GET /me", mw.Auth(http.HandlerFunc(s.mePage)))
	mux.Handle("POST /logout", mw.Auth(http.HandlerFunc(s.logout)))
	return stack(mux)
}

type Server struct {
	DB      *db.SQL
	Indexer *source.Indexer
	env     *models.EnvVars
}

func StartServer() {
	env := models.GetEnv()
	// create db
	db, err := db.GetDB(env)
	if err != nil {
		slog.Error("", "err", err)
		panic(err)
	}
	idxr := source.NewIndexer(env, db)
	defer db.Close()
	server := &Server{
		DB:      db,
		env:     env,
		Indexer: idxr,
	}

	// create mux
	// start the server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%v:%v", env.Host, env.Port),
		Handler: server.RegisterRoutes(),
	}

	slog.Info("server started", "host", env.Host, "port", env.Port)

	go func() {
		if err := server.Indexer.IndexShows(); err != nil {
			slog.Error("failed to index shows", "err", err)
		}
	}()

	cleanUpTicker := time.NewTicker(30 * time.Minute)
	defer cleanUpTicker.Stop()

	go func() {
		for {
			select {
			case <-cleanUpTicker.C:
				err := db.DeleteExpiredCookies()
				if err != nil {
					slog.Error("failed to delete expired cookies", "err", err)
				} else {
					slog.Info("expired cookies deleted")
				}
			}
		}
	}()

	httpServer.ListenAndServe()
}

func (s *Server) render(w http.ResponseWriter, r *http.Request, status int, page string, data ...any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		if len(data) > 0 {
			if err := web.RenderContent(w, page, data[0]); err != nil {
				slog.Error("rendering page fragment failed", "page", page, "err", err)
			}
		} else {
			if err := web.RenderContent(w, page, nil); err != nil {
				slog.Error("rendering page fragment failed", "page", page, "err", err)
			}
		}
		return
	}
	w.WriteHeader(status)
	var pageData any
	if len(data) > 0 {
		pageData = data[0]
	}
	if err := web.RenderBase(w, page, pageData); err != nil {
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

func (s *Server) RenderPageHandler(page string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.render(w, r, http.StatusOK, page)
	}
}
