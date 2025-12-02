package handler

import (
	"itk/internal/service"

	"github.com/go-chi/chi"
)

type Handler struct {
	walletService service.WalletService
}

func NewHandler(walletService service.WalletService) *Handler {
	return &Handler{walletService: walletService}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/wallet", h.handleWalletOperation)
		r.Get("/wallets/{walletId}", h.handleGetBalance)
	})
}
