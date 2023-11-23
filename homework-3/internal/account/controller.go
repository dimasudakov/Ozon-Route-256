package account

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/apperr"
	_ "gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/subscription"
	"net/http"
	"strconv"
)

type Service interface {
	CreateBankAccount(account BankAccount) (*BankAccount, error)
	GetBankAccountById(id int) (*BankAccount, error)
	UpdateBankAccount(id int, account BankAccount) (*BankAccount, error)
	DeleteBankAccount(id int) (*BankAccount, error)
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
	var bankAccountDto BankAccountDto

	if err := json.NewDecoder(r.Body).Decode(&bankAccountDto); err != nil {
		return apperr.NewBadRequestError("Invalid request body")
	}

	response, err := c.service.CreateBankAccount(bankAccountDto.MapToModel())
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
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	response, err := c.service.GetBankAccountById(id)
	if err != nil {
		return err
	}

	var bankAccountDto BankAccountDto
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
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var bankAccountDto BankAccountDto

	if err := json.NewDecoder(r.Body).Decode(&bankAccountDto); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	response, err := c.service.UpdateBankAccount(id, bankAccountDto.MapToModel())
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
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	response, err := c.service.DeleteBankAccount(id)
	if err != nil {
		return err
	}

	var bankAccountDto BankAccountDto
	jsonResponse, err2 := json.Marshal(bankAccountDto.MapFromModel(*response))
	if err2 != nil {
		return apperr.NewInternalServerError("Internal server error")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

	return nil
}
