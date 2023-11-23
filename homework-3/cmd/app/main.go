package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/account"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/app"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/subscription"
	"log"
	"net/http"
)

func main() {
	config, err := app.InitConfig()
	if err != nil {
		_ = fmt.Errorf("apperr occurred during config initialization: %s", err)
		return
	}

	db, err := app.InitDB(config)
	if err != nil {
		_ = fmt.Errorf("apperr occured during connection to db: %s", err)
		return
	}
	defer db.DB.Close()

	db.UpMigrations()

	bankAccountRepository := account.NewBankAccountRepository(db.DB)
	bankAccountService := account.NewBankAccountService(bankAccountRepository)
	bankAccountController := account.NewBankAccountController(bankAccountService)

	subscriptionRepository := subscription.NewSubscriptionRepository(db.DB)
	subscriptionService := subscription.NewSubscriptionService(subscriptionRepository)
	subscriptionController := subscription.NewSubscriptionController(subscriptionService)

	coreController := app.NewCoreController(bankAccountController, subscriptionController)

	r := mux.NewRouter()

	app.ConfigureRoutes(r, coreController)

	http.Handle("/", r)

	s := app.New()
	if err := s.Run(config.Server.Port); err != nil {
		log.Fatalf("Error occured while running http server: %s", err.Error())
	}

}
