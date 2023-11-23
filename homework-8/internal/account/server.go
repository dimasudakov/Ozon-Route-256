package account

import (
	"context"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/account/model"
	logg "gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/app/logging"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/pkg/bank_accounts"
	"go.uber.org/zap"
)

type BankAccountGrpcImpl struct {
	service Service
	bank_accounts.UnimplementedBankAccountServiceServer
}

func NewBankAccountGrpcImpl(service Service) *BankAccountGrpcImpl {
	return &BankAccountGrpcImpl{
		service: service,
	}
}

func (b BankAccountGrpcImpl) CreateBankAccount(ctx context.Context, request *bank_accounts.CreateBankAccountRequest) (*bank_accounts.CreateBankAccountResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateBankAccount")
	defer span.Finish()

	logger := logg.FromContext(ctx)
	logger.With(
		zap.String("method", "create bank account"),
		zap.Any("request", request),
	)
	ctx = logg.ToContext(ctx, logger)

	accountRequest, err := model.MapFromDto(request.GetAccount())
	if err != nil {
		logg.Errorf(ctx, err.Error())
		return nil, err
	}

	bankAccount, err := b.service.CreateBankAccount(ctx, accountRequest)
	if err != nil {
		logg.Errorf(ctx, err.Error())
		return nil, err
	}

	return &bank_accounts.CreateBankAccountResponse{
		Account: bankAccount.MapToDto(),
	}, nil

}

func (b BankAccountGrpcImpl) GetBankAccountById(ctx context.Context, request *bank_accounts.GetBankAccountByIdRequest) (*bank_accounts.GetBankAccountResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetBankAccountById")
	defer span.Finish()

	logger := logg.FromContext(ctx)
	logger.With(
		zap.String("method", "get bank account by id"),
		zap.Any("request", request),
	)
	ctx = logg.ToContext(ctx, logger)

	id, err := uuid.Parse(request.GetId().GetValue())
	if err != nil {
		logg.Errorf(ctx, err.Error())
		return nil, err
	}

	account, err := b.service.GetBankAccountById(ctx, id)
	if err != nil {
		logg.Errorf(ctx, err.Error())
		return nil, err
	}

	return &bank_accounts.GetBankAccountResponse{
		Account: account.MapToDto(),
	}, nil
}

func (b BankAccountGrpcImpl) UpdateBankAccount(ctx context.Context, request *bank_accounts.UpdateBankAccountRequest) (*bank_accounts.UpdateBankAccountResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UpdateBankAccount")
	defer span.Finish()

	logger := logg.FromContext(ctx)
	logger.With(
		zap.String("method", "update bank account"),
		zap.Any("request", request),
	)
	ctx = logg.ToContext(ctx, logger)

	id, err := uuid.Parse(request.GetId().GetValue())
	if err != nil {
		logg.Errorf(ctx, err.Error())
		return nil, err
	}
	accountRequest, err := model.MapFromDto(request.GetAccount())
	if err != nil {
		logg.Errorf(ctx, err.Error())
		return nil, err
	}

	updatedAccount, err := b.service.UpdateBankAccount(ctx, id, accountRequest)
	if err != nil {
		logg.Errorf(ctx, err.Error())
		return nil, err
	}

	return &bank_accounts.UpdateBankAccountResponse{
		Account: updatedAccount.MapToDto(),
	}, nil
}

func (b BankAccountGrpcImpl) DeleteBankAccount(ctx context.Context, request *bank_accounts.DeleteBankAccountRequest) (*bank_accounts.DeleteBankAccountResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DeleteBankAccount")
	defer span.Finish()

	logger := logg.FromContext(ctx)
	logger.With(
		zap.String("method", "delete bank account"),
		zap.Any("request", request),
	)
	ctx = logg.ToContext(ctx, logger)

	id, err := uuid.Parse(request.GetId().GetValue())
	if err != nil {
		logg.Errorf(ctx, err.Error())
		return nil, err
	}

	deletedAccount, err := b.service.DeleteBankAccount(ctx, id)
	if err != nil {
		logg.Errorf(ctx, err.Error())
		return nil, err
	}

	return &bank_accounts.DeleteBankAccountResponse{
		Account: deletedAccount.MapToDto(),
	}, nil
}
