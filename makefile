.PHONY: build up-agent up-dashboard local-db unit-test integration-test lint-test test

LINT := $(shell command -v golangci-lint 2> /dev/null)

build:
	go build -o bin/agent ./cmd/agent
	go build -o bin/dashboard ./cmd/backend
	go build -o bin/smithy ./cmd/smithy

up-agent:
	go build -o bin/agent ./cmd/agent
	PORT=3000 bin/agent

up-dashboard:
	go build -o bin/dashboard ./cmd/backend
	PORT=2999 bin/dashboard

local-db:
	@docker-compose down
	@docker-compose up -d

integration-test:
	go test ./... -tags=integration -count=1

unit-test:
	go test ./... -tags=unit -count=1

lint-test:
ifndef LINT
		go install ./vendor/github.com/golangci/golangci-lint/cmd/golangci-lint
endif
		golangci-lint run

test: lint-test unit-test integration-test
