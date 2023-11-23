//go:generate mockgen -source=./controller.go -destination=./mocks/service.go -package=mock_account

package account

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/account/dtos"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/account/model"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/apperr"
	_ "gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/subscription"
	"net/http"
)

type Service interface {
	CreateBankAccount(account *model.BankAccount) (*model.BankAccount, error)
	GetBankAccountById(id uuid.UUID) (*model.BankAccount, error)
	UpdateBankAccount(id uuid.UUID, account *model.BankAccount) (*model.BankAccount, error)
	DeleteBankAccount(id uuid.UUID) (*model.BankAccount, error)
}

type BankAccountController struct {
	service Service
}

func NewBankAccountController(service Service) *BankAccountController {
	return &BankAccountController{
		service: service,
	}
}

func (c *BankAccountController) CreateBankAccount(w http.ResponseWriter, r *http.Request) error {
	var bankAccountDto dtos.BankAccountDto

	if err := json.NewDecoder(r.Body).Decode(&bankAccountDto); err != nil {
		return apperr.NewBadRequestError("Invalid request body")
	}

	account := bankAccountDto.MapToModel()
	response, err := c.service.CreateBankAccount(&account)
	if err != nil {
		return err
	}

	jsonResponse, err2 := json.Marshal(bankAccountDto.MapFromModel(*response))
	if err2 != nil {
		return apperr.NewInternalServerError("Internal server error")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

	return nil
}

func (c *BankAccountController) GetBankAccount(w http.ResponseWriter, r *http.Request) error {
	idStr, ok := mux.Vars(r)["id"]
	if !ok {
		return apperr.NewBadRequestError("Bad request")
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return apperr.NewBadRequestError("Bad request")
	}
	response, err := c.service.GetBankAccountById(id)
	if err != nil {
		return err
	}

	var bankAccountDto dtos.BankAccountDto
	jsonResponse, err2 := json.Marshal(bankAccountDto.MapFromModel(*response))
	if err2 != nil {
		return apperr.NewInternalServerError("Internal server error")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

	return nil
}

func (c *BankAccountController) UpdateBankAccount(w http.ResponseWriter, r *http.Request) error {
	idStr, ok := mux.Vars(r)["id"]
	if !ok {
		return apperr.NewBadRequestError("Bad request")
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return apperr.NewBadRequestError("Bad request")
	}

	var bankAccountDto dtos.BankAccountDto

	if err := json.NewDecoder(r.Body).Decode(&bankAccountDto); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	account := bankAccountDto.MapToModel()
	response, err := c.service.UpdateBankAccount(id, &account)
	if err != nil {
		return err
	}

	jsonResponse, err2 := json.Marshal(bankAccountDto.MapFromModel(*response))
	if err2 != nil {
		return apperr.NewInternalServerError("Internal server error")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

	return nil
}

func (c *BankAccountController) DeleteBankAccount(w http.ResponseWriter, r *http.Request) error {
	idStr, ok := mux.Vars(r)["id"]
	if !ok {
		return apperr.NewBadRequestError("Bad request")
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return apperr.NewBadRequestError("Bad request")
	}

	response, err := c.service.DeleteBankAccount(id)
	if err != nil {
		return err
	}

	var bankAccountDto dtos.BankAccountDto
	jsonResponse, err2 := json.Marshal(bankAccountDto.MapFromModel(*response))
	if err2 != nil {
		return apperr.NewInternalServerError("Internal server error")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

	return nil
}
