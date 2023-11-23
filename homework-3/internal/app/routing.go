package app

import (
	"github.com/gorilla/mux"
)

func ConfigureRoutes(r *mux.Router, controller *Controller) {
	configureBankAccountRoutes(r, controller.accountController)
	configureSubscriptionRoutes(r, controller.subscriptionController)
}

func configureBankAccountRoutes(r *mux.Router, accountController AccountController) {
	r.HandleFunc("/bank-accounts", ErrorHandler(accountController.CreateBankAccount)).Methods("POST")
	r.HandleFunc("/bank-accounts/{id:[0-9]+}", ErrorHandler(accountController.GetBankAccount)).Methods("GET")
	r.HandleFunc("/bank-accounts/{id:[0-9]+}", ErrorHandler(accountController.UpdateBankAccount)).Methods("PUT")
	r.HandleFunc("/bank-accounts/{id:[0-9]+}", ErrorHandler(accountController.DeleteBankAccount)).Methods("DELETE")
}

func configureSubscriptionRoutes(r *mux.Router, subscriptionController SubscriptionController) {
	r.HandleFunc("/subscriptions", ErrorHandler(subscriptionController.CreateSubscription)).Methods("POST")
	r.HandleFunc("/subscriptions/{id:[0-9]+}", ErrorHandler(subscriptionController.GetSubscription)).Methods("GET")
	r.HandleFunc("/subscriptions", ErrorHandler(subscriptionController.UpdateSubscription)).Methods("PUT")
	r.HandleFunc("/subscriptions/{id:[0-9]+}", ErrorHandler(subscriptionController.DeleteSubscription)).Methods("DELETE")
}
