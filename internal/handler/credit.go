package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bank-api/internal/middleware"
	"bank-api/internal/service"

	"github.com/gorilla/mux"
)

type CreditHandler struct {
	creditService *service.CreditService
}

func NewCreditHandler(svc *service.CreditService) *CreditHandler {
	return &CreditHandler{creditService: svc}
}

func (h *CreditHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	var req struct {
		AccountID  int     `json:"account_id"`
		Amount     float64 `json:"amount"`
		TermMonths int     `json:"term_months"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	credit, err := h.creditService.CreateCredit(r.Context(), userID, req.AccountID, req.Amount, req.TermMonths)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(credit)
}

func (h *CreditHandler) GetSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	creditID, err := strconv.Atoi(vars["creditId"])
	if err != nil {
		http.Error(w, "invalid creditId", http.StatusBadRequest)
		return
	}

	schedule, err := h.creditService.GetSchedule(r.Context(), creditID)
	if err != nil {
		http.Error(w, "failed to get schedule", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schedule) // даже если schedule пуст, будет []
}
