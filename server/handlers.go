package server

import "net/http"

func (s *Server) RegisterRoutes() http.Handler {

	mux := http.NewServeMux()
	mw := &MiddleWares{
		DB: s.DB,
	}
	NoAuthStack := CreateStack(
		mw.Logging,
		mw.Recovery,
	)

	// register handlers
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	handler := NoAuthStack(mux)
	return handler

}
