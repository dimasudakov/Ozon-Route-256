//go:generate mockgen  -destination=../account/mocks/controller.go -package=mock_account -source=./controller.go AccountController

package app

import "net/http"

type Controller struct {
	AccountController      AccountController
	SubscriptionController SubscriptionController
}

type AccountController interface {
	CreateBankAccount(w http.ResponseWriter, r *http.Request) error
	GetBankAccount(w http.ResponseWriter, r *http.Request) error
	UpdateBankAccount(w http.ResponseWriter, r *http.Request) error
	DeleteBankAccount(w http.ResponseWriter, r *http.Request) error
}

type SubscriptionController interface {
	CreateSubscription(w http.ResponseWriter, r *http.Request) error
	GetSubscription(w http.ResponseWriter, r *http.Request) error
	UpdateSubscription(w http.ResponseWriter, r *http.Request) error
	DeleteSubscription(w http.ResponseWriter, r *http.Request) error
}

func NewCoreController(accountController AccountController, subscriptionController SubscriptionController) *Controller {
	return &Controller{
		AccountController:      accountController,
		SubscriptionController: subscriptionController,
	}
}
