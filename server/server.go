package server

import (
	"net/http"

	"github.com/fintech-asean/fintech-go/account"
	kitlog "github.com/go-kit/kit/log"

	"github.com/go-chi/chi"
)

// Server holds the dependencies for a HTTP server.
type Server struct {
	Account account.Service
	Logger  kitlog.Logger

	router chi.Router
}

// New returns a new HTTP server.
func New(bs account.Service, logger kitlog.Logger) *Server {
	s := &Server{}

	r := chi.NewRouter()

	r.Use(accessControl)

	r.Route("/account", func(r chi.Router) {
		h := accountHandler{s.Account, s.Logger}
		r.Mount("/v1", h.router())
	})

	s.router = r
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
