THIS_FILE := $(lastword $(MAKEFILE_LIST))
DOCKER_COMPOSE_FILE=$(CURDIR)/docker/docker-compose.yaml
POSTGRES_CONNECTION = "user=postgres password=postgres dbname=homework-5 host=localhost port=6666 sslmode=disable"
MIGRATIONS_DIR := db/migrations


docker-up:
	@docker-compose -f "$(DOCKER_COMPOSE_FILE)" up

docker-up-kafka:
	@docker-compose -f "$(DOCKER_COMPOSE_FILE)" up -d zookeeper kafka1 kafka2

docker-down:
	@docker-compose -f "$(DOCKER_COMPOSE_FILE)" down

docker-start:
	@docker-compose -f "$(DOCKER_COMPOSE_FILE)" start

docker-stop:
	@docker-compose -f "$(DOCKER_COMPOSE_FILE)" stop

docker-ps:
	@docker-compose -f "$(DOCKER_COMPOSE_FILE)" ps

docker-restart:
	@docker-compose -f "$(DOCKER_COMPOSE_FILE)" restart


migrate-up:
	@goose -dir $(MIGRATIONS_DIR) postgres $(POSTGRES_CONNECTION) up

migrate-down:
	@goose -dir $(MIGRATIONS_DIR) postgres $(POSTGRES_CONNECTION) down

unit-test: docker-start
	@go test ./internal/... -cover

integration-test: docker-start
	@go test -tags=integration ./tests/... -coverpkg=./internal/...





run: docker-start
	@go run cmd/app/main.go


.PHONY: docker-up docker-down docker-start docker-stop docker-ps docker-restart
.PHONY: migrate-up migrate-down create-migration migrate-up-test-db
.PHONY: unit-test
