package subscription

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/pkg/subscriptions"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Subscription struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"subscription_name"`
	Price     int       `db:"price"`
	StartDate time.Time `db:"start_date"`
	EndDate   time.Time `db:"end_date"`
	AccountID uuid.UUID `db:"account_id"`
}

func (s Subscription) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name, validation.Required),
		validation.Field(&s.Price, validation.Required, validation.Min(0)),
		validation.Field(&s.EndDate, validation.Min(time.Now())),
	)
}

func MapFromDto(dto *subscriptions.SubscriptionDto) (*Subscription, error) {
	id, err := uuid.Parse(dto.GetId().GetValue())
	if err != nil {
		return nil, err
	}
	accountId, _ := uuid.Parse(dto.GetAccountId().GetValue())
	return &Subscription{
		ID:        id,
		Name:      dto.GetSubscriptionName(),
		Price:     int(dto.GetPrice()),
		StartDate: dto.GetStartDate().GetValue().AsTime(),
		EndDate:   dto.GetEndDate().GetValue().AsTime(),
		AccountID: accountId,
	}, nil
}

func (s Subscription) MapToDto() *subscriptions.SubscriptionDto {

	return &subscriptions.SubscriptionDto{
		Id:               &subscriptions.UUID{Value: s.ID.String()},
		SubscriptionName: s.Name,
		Price:            int32(s.Price),
		StartDate:        &subscriptions.Timestamp{Value: timestamppb.New(s.StartDate)},
		EndDate:          &subscriptions.Timestamp{Value: timestamppb.New(s.EndDate)},
		AccountId:        &subscriptions.UUID{Value: s.AccountID.String()},
	}
}

func MapFromDtoList(dto []*subscriptions.SubscriptionDto) ([]Subscription, error) {
	result := make([]Subscription, len(dto))
	for i, dtoValue := range dto {
		sub, err := MapFromDto(dtoValue)
		if err != nil {
			return nil, err
		}
		result[i] = *sub
	}
	return result, nil
}

func MapToDtoList(subs []Subscription) []*subscriptions.SubscriptionDto {
	result := make([]*subscriptions.SubscriptionDto, len(subs))
	for i, sub := range subs {
		result[i] = sub.MapToDto()
	}
	return result
}
