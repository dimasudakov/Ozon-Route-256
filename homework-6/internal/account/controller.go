//go:generate mockgen -source=./controller.go -destination=./mocks/service.go -package=mock_account

package account

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/account/dtos"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/account/model"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/app/logging"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/apperr"
	_ "gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/subscription"
	"io"
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
	logger  logging.Logger
}

func NewBankAccountController(service Service, logger logging.Logger) *BankAccountController {
	return &BankAccountController{
		service: service,
		logger:  logger,
	}
}

func (c *BankAccountController) CreateBankAccount(w http.ResponseWriter, r *http.Request) error {
	body, _ := io.ReadAll(r.Body)
	logMsg := logging.LogMessage{
		RequestURI:  r.RequestURI,
		RequestType: r.Method,
		Method:      "CreateBankAccount",
		Body:        string(body),
	}
	c.logger.Log(logMsg)

	var bankAccountDto dtos.BankAccountDto

	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&bankAccountDto); err != nil {
		logMsg.Info = "Bad request"
		c.logger.Error(logMsg)
		return apperr.NewBadRequestError("Invalid request body")
	}

	account := bankAccountDto.MapToModel()
	response, err := c.service.CreateBankAccount(&account)
	if err != nil {
		logMsg.Info = err.Error()
		c.logger.Error(logMsg)
		return err
	}

	jsonResponse, err2 := json.Marshal(bankAccountDto.MapFromModel(*response))
	if err2 != nil {
		logMsg.Info = err2.Error()
		c.logger.Error(logMsg)
		return apperr.NewInternalServerError("Internal server error")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

	return nil
}

func (c *BankAccountController) GetBankAccount(w http.ResponseWriter, r *http.Request) error {
	logMsg := logging.LogMessage{
		RequestURI:  r.RequestURI,
		RequestType: r.Method,
		Method:      "GetBankAccount",
	}
	c.logger.Log(logMsg)

	idStr, ok := mux.Vars(r)["id"]
	if !ok {
		logMsg.Info = "Bad request"
		c.logger.Error(logMsg)
		return apperr.NewBadRequestError("Bad request")
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		logMsg.Info = "Bad request"
		c.logger.Error(logMsg)
		return apperr.NewBadRequestError("Bad request")
	}

	response, err := c.service.GetBankAccountById(id)
	if err != nil {
		logMsg.Info = err.Error()
		c.logger.Error(logMsg)
		return err
	}

	var bankAccountDto dtos.BankAccountDto
	jsonResponse, err2 := json.Marshal(bankAccountDto.MapFromModel(*response))
	if err2 != nil {
		logMsg.Info = err2.Error()
		c.logger.Error(logMsg)
		return apperr.NewInternalServerError("Internal server error")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

	return nil
}

func (c *BankAccountController) UpdateBankAccount(w http.ResponseWriter, r *http.Request) error {
	body, _ := io.ReadAll(r.Body)
	logMsg := logging.LogMessage{
		RequestURI:  r.RequestURI,
		RequestType: r.Method,
		Method:      "UpdateBankAccount",
		Body:        string(body),
	}
	c.logger.Log(logMsg)

	idStr, ok := mux.Vars(r)["id"]
	if !ok {
		logMsg.Info = "Bad request"
		c.logger.Error(logMsg)
		return apperr.NewBadRequestError("Bad request")
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		logMsg.Info = err.Error()
		c.logger.Error(logMsg)
		return apperr.NewBadRequestError("Bad request")
	}

	var bankAccountDto dtos.BankAccountDto

	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&bankAccountDto); err != nil {
		logMsg.Info = "Invalid request body"
		c.logger.Error(logMsg)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	account := bankAccountDto.MapToModel()
	response, err := c.service.UpdateBankAccount(id, &account)
	if err != nil {
		logMsg.Info = err.Error()
		c.logger.Error(logMsg)
		return err
	}

	jsonResponse, err2 := json.Marshal(bankAccountDto.MapFromModel(*response))
	if err2 != nil {
		logMsg.Info = err2.Error()
		c.logger.Error(logMsg)
		return apperr.NewInternalServerError("Internal server error")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

	return nil
}

func (c *BankAccountController) DeleteBankAccount(w http.ResponseWriter, r *http.Request) error {
	logMsg := logging.LogMessage{
		RequestURI:  r.RequestURI,
		RequestType: r.Method,
		Method:      "DeleteBankAccount",
	}
	c.logger.Log(logMsg)

	idStr, ok := mux.Vars(r)["id"]
	if !ok {
		logMsg.Info = "Bad request"
		c.logger.Error(logMsg)
		return apperr.NewBadRequestError("Bad request")
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		logMsg.Info = err.Error()
		c.logger.Error(logMsg)
		return apperr.NewBadRequestError("Bad request")
	}

	response, err := c.service.DeleteBankAccount(id)
	if err != nil {
		logMsg.Info = err.Error()
		c.logger.Error(logMsg)
		return err
	}

	var bankAccountDto dtos.BankAccountDto
	jsonResponse, err2 := json.Marshal(bankAccountDto.MapFromModel(*response))
	if err2 != nil {
		logMsg.Info = err2.Error()
		c.logger.Error(logMsg)
		return apperr.NewInternalServerError("Internal server error")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

	return nil
}
