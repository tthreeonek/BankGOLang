package service

import (
	"context"
	"time"
)

type AnalyticsService struct {
	// Здесь могут быть зависимости, например, репозиторий транзакций
}

func NewAnalyticsService() *AnalyticsService {
	return &AnalyticsService{}
}

// MonthlyStats возвращает статистику доходов/расходов за указанный месяц
func (s *AnalyticsService) MonthlyStats(ctx context.Context, userID int, month time.Time) (map[string]float64, error) {
	// TODO: заменить на реальный запрос из БД
	stats := map[string]float64{
		"income":  50000.0,
		"expense": 30000.0,
		"profit":  20000.0,
	}
	return stats, nil
}

// CreditLoad возвращает аналитику по кредитной нагрузке
func (s *AnalyticsService) CreditLoad(ctx context.Context, userID int) (map[string]interface{}, error) {
	// Заглушка
	return map[string]interface{}{
		"active_credits": 1,
		"total_debt":     150000.0,
	}, nil
}
