package server

import (
	"backend/db"
	"backend/utils"
	"database/sql"
	"log/slog"
	"net/http"
	"time"
)

type MiddleWares struct {
	DB *db.SQL
}

type Middleware func(http.Handler) http.Handler

func CreateStack(handlers ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(handlers) - 1; i >= 0; i-- {
			handler := handlers[i]
			next = handler(next)
		}
		return next
	}
}

func (mw *MiddleWares) Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			r := recover()
			if r == nil {
				return
			}
			if err, isSQLError := r.(error); isSQLError && err == sql.ErrNoRows {
				slog.Error("error", "data", r)
				return
			}
			// log the error and return a 500
			slog.Error("error", "data", r)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}()
		next.ServeHTTP(w, r)
	})
}
func (mw *MiddleWares) cacheStaticFiles(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		next.ServeHTTP(w, r)
	})
}

func (mw *MiddleWares) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := utils.GetUserByCookie(r, mw.DB)
		if err != nil {
			slog.Error("error happened", "err", err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (ww *wrappedWriter) WriteHeader(code int) {
	ww.statusCode = code
	ww.ResponseWriter.WriteHeader(code)
}

func (mw *MiddleWares) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		if r.URL.Path == "/posts" {
			next.ServeHTTP(w, r)
			return
		}

		startTime := time.Now()

		next.ServeHTTP(wrapped, r)

		slog.Info("",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Duration("duration", time.Since(startTime)),
			slog.Int("status", wrapped.statusCode),
		)
	})
}
