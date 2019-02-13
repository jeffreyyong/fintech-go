package rest

import (
	"net/http"

	"github.com/go-chi/chi"
	kitlog "github.com/go-kit/kit/log"

	"github.com/jeffreyyong/fintech-go/lib/msgqueue"
	"github.com/jeffreyyong/fintech-go/lib/persistence"
)

// Server holds the dependencies for a HTTP server.
type Server struct {
	Logger kitlog.Logger
	router chi.Router
}

// New returns a new HTTP server for accounts.
func New(dbHandler persistence.DatabaseHandler, eventEmitter msgqueue.EventEmitter, logger kitlog.Logger) *Server {
	s := &Server{
		Logger: logger,
	}

	handler := newAccountHandler(dbHandler, eventEmitter)

	r := chi.NewRouter()
	r.Use(accessControl)

	r.Get("/accounts", handler.allAccountHandler)
	r.Post("/account", handler.newAccountHandler)
	r.Get("/account/{accountID}", handler.oneAccountHandler)

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
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
