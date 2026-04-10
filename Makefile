SHELL := /bin/bash

GO ?= go
COMPOSE := $(shell docker compose version >/dev/null 2>&1 && echo "docker compose" || echo "docker-compose")
GO_ENV := GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod
LINT_ENV := $(GO_ENV) GOLANGCI_LINT_CACHE=/tmp/golangci-lint-cache

GOOSE ?= goose
DB_URL ?= postgres://postgres:postgres@localhost:5432/wishlist_db?sslmode=disable

.PHONY: run up down logs test test-cover fmt lint lint-fix tidy migrate-up migrate-down migrate-status

run:
	$(GO_ENV) $(GO) run ./cmd/wishlist-service

up:
	$(COMPOSE) up --build

down:
	$(COMPOSE) down -v

logs:
	$(COMPOSE) logs -f app

test:
	$(GO_ENV) $(GO) test ./...

test-cover:
	$(GO_ENV) $(GO) test -cover ./...

fmt:
	$(GO) fmt ./...

lint:
	$(LINT_ENV) golangci-lint run ./...

lint-fix:
	$(LINT_ENV) golangci-lint run --fix ./...

tidy:
	$(GO_ENV) $(GO) mod tidy

migrate-up:
	$(GOOSE) -dir ./migrations postgres "$(DB_URL)" up

migrate-down:
	$(GOOSE) -dir ./migrations postgres "$(DB_URL)" down

migrate-status:
	$(GOOSE) -dir ./migrations postgres "$(DB_URL)" status
