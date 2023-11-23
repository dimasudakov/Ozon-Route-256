package subscription

import (
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/apperr"
	"time"
)

type Repository interface {
	CreateSubscription(subscription Subscription) (*Subscription, error)
	GetSubscriptionByID(id int) (*Subscription, error)
}

type SubscriptionService struct {
	repository Repository
}

func NewSubscriptionService(repository Repository) *SubscriptionService {
	return &SubscriptionService{
		repository: repository,
	}
}

func (s *SubscriptionService) CreateSubscription(subscription Subscription) (*Subscription, error) {
	if err := subscription.Validate(); err != nil {
		return nil, apperr.NewBadRequestError(err.Error())
	}
	subscription.StartDate = time.Now()

	createdSubscription, err := s.repository.CreateSubscription(subscription)
	if err != nil {
		return nil, err
	}

	return createdSubscription, nil
}

func (s *SubscriptionService) GetSubscriptionById(id int) (*Subscription, error) {
	subscription, err := s.repository.GetSubscriptionByID(id)
	if err != nil {
		return nil, err
	}
	return subscription, nil
}
