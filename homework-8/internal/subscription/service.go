package subscription

import (
	"context"
	"github.com/google/uuid"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/apperr"
	"time"
)

type Service interface {
	CreateSubscription(ctx context.Context, subscription Subscription) (*Subscription, error)
	GetSubscriptionById(ctx context.Context, id uuid.UUID) (*Subscription, error)
}

type SubscriptionService struct {
	repository Repository
}

func NewSubscriptionService(repository Repository) *SubscriptionService {
	return &SubscriptionService{
		repository: repository,
	}
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, subscription Subscription) (*Subscription, error) {
	if err := subscription.Validate(); err != nil {
		return nil, apperr.NewBadRequestError(err.Error())
	}
	subscription.StartDate = time.Now()

	createdSubscription, err := s.repository.CreateSubscription(ctx, subscription)
	if err != nil {
		return nil, err
	}

	return createdSubscription, nil
}

func (s *SubscriptionService) GetSubscriptionById(ctx context.Context, id uuid.UUID) (*Subscription, error) {
	subscription, err := s.repository.GetSubscriptionByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return subscription, nil
}
