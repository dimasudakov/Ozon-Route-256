package subscription

import (
	"context"
	"github.com/google/uuid"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/pkg/subscriptions"
)

type SubscriptionGrpcImpl struct {
	service Service
	subscriptions.UnimplementedSubscriptionServiceServer
}

func NewSubscriptionGrpcImpl(service Service) *SubscriptionGrpcImpl {
	return &SubscriptionGrpcImpl{
		service: service,
	}
}

func (s SubscriptionGrpcImpl) CreateSubscription(ctx context.Context, request *subscriptions.CreateSubscriptionRequest) (*subscriptions.CreateSubscriptionResponse, error) {
	subscriptionRequest, err := MapFromDto(request.GetSubscription())
	if err != nil {
		return nil, err
	}

	createdSubscription, err := s.service.CreateSubscription(ctx, *subscriptionRequest)
	if err != nil {
		return nil, err
	}

	return &subscriptions.CreateSubscriptionResponse{
		Subscription: createdSubscription.MapToDto(),
	}, nil
}

func (s SubscriptionGrpcImpl) GetSubscriptionById(ctx context.Context, request *subscriptions.GetSubscriptionByIdRequest) (*subscriptions.GetSubscriptionResponse, error) {
	id, err := uuid.Parse(request.GetId().GetValue())
	if err != nil {
		return nil, err
	}

	subscription, err := s.service.GetSubscriptionById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &subscriptions.GetSubscriptionResponse{
		Subscription: subscription.MapToDto(),
	}, nil
}
