package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"bank-api/internal/middleware"
	"bank-api/internal/service"
)

type AnalyticsHandler struct {
	analyticsService *service.AnalyticsService
}

func NewAnalyticsHandler(svc *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsService: svc}
}

// MonthlyStats возвращает статистику за указанный месяц.
// В запросе параметры month (YYYY-MM) можно передать через query, либо по умолчанию текущий месяц.
func (h *AnalyticsHandler) MonthlyStats(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	// Простейшая обработка параметра month из query (?month=2025-01)
	monthParam := r.URL.Query().Get("month")
	var month time.Time
	if monthParam != "" {
		parsed, err := time.Parse("2006-01", monthParam)
		if err != nil {
			http.Error(w, "invalid month format, use YYYY-MM", http.StatusBadRequest)
			return
		}
		month = parsed
	} else {
		now := time.Now()
		month = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	}

	stats, err := h.analyticsService.MonthlyStats(r.Context(), userID, month)
	if err != nil {
		http.Error(w, "failed to get analytics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// CreditLoad возвращает информацию о кредитной нагрузке пользователя
func (h *AnalyticsHandler) CreditLoad(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	load, err := h.analyticsService.CreditLoad(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to get credit load", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(load)
}
