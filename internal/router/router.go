package router

import (
	"bank-api/internal/handler"
	"bank-api/internal/middleware"

	"github.com/gorilla/mux"
)

func NewRouter(
	authH *handler.AuthHandler,
	accountH *handler.AccountHandler,
	cardH *handler.CardHandler,
	transferH *handler.TransferHandler,
	creditH *handler.CreditHandler,
	analyticsH *handler.AnalyticsHandler,
	jwtSecret string,
) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/register", authH.Register).Methods("POST")
	r.HandleFunc("/login", authH.Login).Methods("POST")

	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware(jwtSecret))

	api.HandleFunc("/accounts", accountH.Create).Methods("POST")
	api.HandleFunc("/accounts", accountH.List).Methods("GET")
	api.HandleFunc("/accounts/{accountId}/deposit", accountH.Deposit).Methods("POST")
	api.HandleFunc("/accounts/{accountId}/predict", accountH.Predict).Methods("GET") // новый
	api.HandleFunc("/cards", cardH.Create).Methods("POST")
	api.HandleFunc("/transfer", transferH.Handle).Methods("POST")
	api.HandleFunc("/credits", creditH.Create).Methods("POST")
	api.HandleFunc("/credits/{creditId}/schedule", creditH.GetSchedule).Methods("GET")
	api.HandleFunc("/analytics/monthly", analyticsH.MonthlyStats).Methods("GET")
	api.HandleFunc("/analytics/credit-load", analyticsH.CreditLoad).Methods("GET")

	return r
}
