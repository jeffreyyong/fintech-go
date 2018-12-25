package server

import (
	"github.com/fintech-asean/fintech-go/balance"
	"github.com/go-chi/chi"
)

type balanceHandler struct {
	s balance.Service

	logger kitlog.logger
}

func (h *balanceHandler) router() chi.Router {

}
