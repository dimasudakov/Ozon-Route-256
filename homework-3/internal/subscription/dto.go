package subscription

import (
	"time"
)

type SubscriptionDto struct {
	ID        int       `json:"id"`
	Name      string    `json:"subscription_name"`
	Price     int       `json:"price"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	AccountID int       `json:"account_id"`
}

func (s *SubscriptionDto) MapFromModel(model Subscription) *SubscriptionDto {

	s.ID = model.ID
	s.Name = model.Name
	s.Price = model.Price
	s.StartDate = model.StartDate
	s.EndDate = model.EndDate
	s.AccountID = model.AccountID

	return s
}

func (s *SubscriptionDto) MapToModel() Subscription {

	return Subscription{
		ID:        s.ID,
		Name:      s.Name,
		Price:     s.Price,
		StartDate: s.StartDate,
		EndDate:   s.EndDate,
		AccountID: s.AccountID,
	}

}

func MapFromModels(subs []Subscription) []SubscriptionDto {
	dtos := make([]SubscriptionDto, 0, len(subs))
	for _, sub := range subs {
		dto := SubscriptionDto{}
		dto.MapFromModel(sub)
		dtos = append(dtos, dto)
	}
	return dtos
}

func MapToModels(subDtos []SubscriptionDto) []Subscription {
	models := make([]Subscription, 0, len(subDtos))
	for _, sub := range subDtos {
		models = append(models, sub.MapToModel())
	}
	return models
}
