package app

import (
	"github.com/gorilla/mux"
	"net/http"
)

type errorHandlerFunc func(f func(w http.ResponseWriter, r *http.Request) error) func(http.ResponseWriter, *http.Request)

func ConfigureRoutes(r *mux.Router, controller *Controller) {
	errorHandler := ErrorHandler
	configureBankAccountRoutes(r, controller.AccountController, errorHandler)
	configureSubscriptionRoutes(r, controller.SubscriptionController, errorHandler)
}

func configureBankAccountRoutes(r *mux.Router, accountController AccountController, errorHandler errorHandlerFunc) {
	r.HandleFunc("/bank-accounts", errorHandler(accountController.CreateBankAccount)).Methods("POST")
	r.HandleFunc("/bank-accounts/{id:[0-9a-fA-F-]+}", errorHandler(accountController.GetBankAccount)).Methods("GET")
	r.HandleFunc("/bank-accounts/{id:[0-9a-fA-F-]+}", errorHandler(accountController.UpdateBankAccount)).Methods("PUT")
	r.HandleFunc("/bank-accounts/{id:[0-9a-fA-F-]+}", errorHandler(accountController.DeleteBankAccount)).Methods("DELETE")
}

func configureSubscriptionRoutes(r *mux.Router, subscriptionController SubscriptionController, errorHandler errorHandlerFunc) {
	r.HandleFunc("/subscriptions", errorHandler(subscriptionController.CreateSubscription)).Methods("POST")
	r.HandleFunc("/subscriptions/{id:[0-9a-fA-F-]+}", errorHandler(subscriptionController.GetSubscription)).Methods("GET")
	r.HandleFunc("/subscriptions/{id:[0-9a-fA-F-]+}", errorHandler(subscriptionController.UpdateSubscription)).Methods("PUT")
	r.HandleFunc("/subscriptions/{id:[0-9a-fA-F-]+}", errorHandler(subscriptionController.DeleteSubscription)).Methods("DELETE")
}
