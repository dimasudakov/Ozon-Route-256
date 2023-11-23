package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/account"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/app"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/app/database"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/internal/subscription"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/pkg/bank_accounts"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-3/pkg/subscriptions"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

var (
	ctx = context.Background()
)

func main() {
	config, err := app.InitConfig()
	if err != nil {
		fmt.Printf("apperr occurred during config initialization: %s", err)
		return
	}

	db, err := database.InitDB(config)
	if err != nil {
		fmt.Printf("error occured during connection to db: %s", err)
		return
	}
	defer db.Close()

	if err := db.UpMigrations(); err != nil {
		fmt.Printf("error occured while performing database migrations: %s", err)
		return
	}

	bankAccountRepository := account.NewBankAccountRepository(db)
	bankAccountService := account.NewBankAccountService(bankAccountRepository)

	subscriptionRepository := subscription.NewSubscriptionRepository(db)
	subscriptionService := subscription.NewSubscriptionService(subscriptionRepository)

	go func() {
		err := runGatewayServer(ctx, config.Server.GatewayPort)
		if err != nil {
			log.Fatal(err)
		}
	}()

	if err := run(ctx, config.Server.GrpcPort, *bankAccountService, *subscriptionService); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, addr string, bankAccountService account.BankAccountService, subscriptionService subscription.SubscriptionService) error {

	setupTracing()

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			app.ContextPropagationUnaryServerInterceptor(),
			app.UnaryErrorHandlerInterceptor(),
		),
	)

	bank_accounts.RegisterBankAccountServiceServer(grpcServer, account.NewBankAccountGrpcImpl(&bankAccountService))
	subscriptions.RegisterSubscriptionServiceServer(grpcServer, subscription.NewSubscriptionGrpcImpl(&subscriptionService))

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	log.Printf("bankAccountService messages listening on %q", addr)
	return grpcServer.Serve(lis)
}

func runGatewayServer(ctx context.Context, addr string) error {
	conn, err := grpc.DialContext(
		context.Background(),
		"127.0.0.1:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	grpcMux := runtime.NewServeMux()
	err = bank_accounts.RegisterBankAccountServiceHandler(ctx, grpcMux, conn)
	if err != nil {
		log.Fatal(err)
	}

	gwServer := &http.Server{
		Addr:    addr,
		Handler: grpcMux,
	}

	return gwServer.ListenAndServe()
}

func setupTracing() {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            false,
			BufferFlushInterval: 1 * time.Second,
		},
	}
	tracer, closer, err := cfg.New(
		"bank-account-service",
	)
	if err != nil {
		fmt.Printf("cannot create tracer: %v\n", err)
		os.Exit(1)
	}
	defer closer.Close()

	opentracing.SetGlobalTracer(tracer)
}
