package subscription

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/apperr"
	"net/http"
	"strconv"
)

type Service interface {
	CreateSubscription(subscription Subscription) (*Subscription, error)
	GetSubscriptionById(id int) (*Subscription, error)
}

type SubscriptionController struct {
	service Service
}

func NewSubscriptionController(service Service) *SubscriptionController {
	return &SubscriptionController{
		service: service,
	}
}

func (c *SubscriptionController) CreateSubscription(w http.ResponseWriter, r *http.Request) error {
	var subscriptionDto SubscriptionDto

	if err := json.NewDecoder(r.Body).Decode(&subscriptionDto); err != nil {
		return apperr.NewBadRequestError("Invalid request body")
	}

	response, err := c.service.CreateSubscription(subscriptionDto.MapToModel())
	if err != nil {
		return err
	}

	jsonResponse, err2 := json.Marshal(subscriptionDto.MapFromModel(*response))
	if err2 != nil {
		apperr.NewInternalServerError("Internal server error")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

	return nil
}

func (c *SubscriptionController) GetSubscription(w http.ResponseWriter, r *http.Request) error {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	response, err := c.service.GetSubscriptionById(id)
	if err != nil {
		return err
	}

	var subscriptionDto SubscriptionDto
	jsonResponse, _ := json.Marshal(subscriptionDto.MapFromModel(*response))
	if err != nil {
		return apperr.NewInternalServerError("Failed to marshal json")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	return nil
}

func (c *SubscriptionController) UpdateSubscription(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (c *SubscriptionController) DeleteSubscription(w http.ResponseWriter, r *http.Request) error {
	return nil
}
