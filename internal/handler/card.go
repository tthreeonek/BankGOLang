package handler

import (
	"bank-api/internal/middleware"
	"bank-api/internal/service"
	"encoding/json"
	"net/http"
)

type CardHandler struct {
	cardService *service.CardService
}

func NewCardHandler(svc *service.CardService) *CardHandler {
	return &CardHandler{cardService: svc}
}

func (h *CardHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	var req struct {
		AccountID int `json:"account_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	card, err := h.cardService.CreateCard(r.Context(), userID, req.AccountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"masked_pan": card.MaskedPAN})
}
