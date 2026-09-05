package server

import (
	"backend/db"
	"backend/models"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
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

		page, err := buildLibraryPage(s.DB, s.env.TvDir)
		if err != nil {
			slog.Error("failed to build library page", "err", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		s.render(w, r, http.StatusOK, "home.html", map[string]any{
			"Shows": page,
		})
	})

	mux.HandleFunc("GET /api/poster/{id...}", s.GetPoster)
	mux.HandleFunc("GET /shows/stream/{fname...}", s.StreamShow)
	mux.HandleFunc("GET /shows/{id...}", s.GetShow)

	mux.HandleFunc("GET /login", s.RenderPageHandler("login.html"))
	mux.HandleFunc("POST /login", s.login)

	mux.HandleFunc("GET /register", s.RenderPageHandler("register.html"))
	mux.HandleFunc("POST /register", s.register)

	mux.Handle("GET /me", mw.Auth(http.HandlerFunc(s.mePage)))
	mux.Handle("POST /logout", mw.Auth(http.HandlerFunc(s.logout)))
	return stack(mux)
}

func (s *Server) StreamShow(w http.ResponseWriter, r *http.Request) {
	rawFname := r.PathValue("fname")

	fname, err := url.PathUnescape(rawFname)
	if err != nil {
		http.Error(w, "invalid path encoding", http.StatusBadRequest)
		return
	}

	cleanPath := path.Clean(filepath.ToSlash(filepath.Join(s.env.TvDir, fname)))
	file, err := os.Open(cleanPath)
	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil || stat.IsDir() {
		http.Error(w, "invalid file", http.StatusNotFound)
		return
	}

	// 4. Assert ReadSeeker for byte-range support in mpv
	http.ServeContent(w, r, stat.Name(), stat.ModTime(), file)
}

type Server struct {
	DB  *db.SQL
	env *models.EnvVars
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
		DB:  db,
		env: env,
	}

	// create mux
	// start the server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%v:%v", env.Host, env.Port),
		Handler: server.RegisterRoutes(),
	}

	slog.Info("server started", "host", env.Host, "port", env.Port)

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
