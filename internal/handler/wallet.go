package handler

import (
	"encoding/json"
	errs "itk/internal/errors"
	"itk/models"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (h *Handler) handleWalletOperation(w http.ResponseWriter, r *http.Request) {
	var req models.WalletOperationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	balance, err := h.walletService.ProcessOperation(r.Context(), req)
	if err != nil {
		switch err {
		case errs.ErrInvalidOperationType:
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		case errs.ErrInsufficientFunds:
			http.Error(w, err.Error(), http.StatusConflict)
			return
		case errs.ErrWalletNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}

	resp := models.WalletBalanceResponse{
		WalletID: req.WalletID,
		Balance:  balance,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) handleGetBalance(w http.ResponseWriter, r *http.Request) {
	walletIDStr := chi.URLParam(r, "walletId")
	id, err := uuid.Parse(walletIDStr)
	if err != nil {
		http.Error(w, "invalid wallet id", http.StatusBadRequest)
		return
	}

	balance, err := h.walletService.GetBalance(r.Context(), id)
	if err != nil {
		switch err {
		case errs.ErrWalletNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}

	resp := models.WalletBalanceResponse{
		WalletID: id,
		Balance:  balance,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
