.EXPORT_ALL_VARIABLES:

APP_NAME=dependabot
PACKAGE=scp-dependency-manager
CURRENT_DIR=$(shell pwd)

COVERAGE_DIR=./coverage

DOCKER_COMPOSE_FILE=$(CURRENT_DIR)/test/docker-compose.yml

VERSION=$(shell cat ${CURRENT_DIR}/VERSION)
BUILD_DATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT=$(shell git rev-parse --short HEAD)
GIT_TAG=$(shell if [ -z "`git status --porcelain`" ]; then git describe --exact-match --tags HEAD 2>/dev/null; fi)

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

$(shell mkdir -p $(COVERAGE_DIR))

.PHONY: test
test:
	go test -v -cover -race -coverprofile=$(COVERAGE_DIR)/coverage.out -covermode=atomic ./...
	go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	gocover-cobertura < $(COVERAGE_DIR)/coverage.out > $(COVERAGE_DIR)/coverage.xml

.PHONY: itest
itest:
	go test --tags=integration -race -coverprofile=$(COVERAGE_DIR)/coverage.out -covermode=atomic -p 1 -cover ./...
	go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	gocover-cobertura < $(COVERAGE_DIR)/coverage.out > $(COVERAGE_DIR)/coverage.xml

.PHONY: fmt
fmt:
	go fmt ./...
	goimports -format-only -w ./..

.PHONY: vendor
vendor:
	go mod vendor

.PHONY: build.binaries
build.binaries:
	goreleaser release --skip-publish --rm-dist --snapshot

# Docker Compose Integration Test tasks
.PHONY:  compose-up
compose-up:
	docker-compose -f $(DOCKER_COMPOSE_FILE) -p $(APP_NAME) up -d

.PHONY:  compose-stop
compose-stop:
	docker-compose -f $(DOCKER_COMPOSE_FILE) -p $(APP_NAME) stop

.PHONY:  compose-destroy
compose-destroy:
	docker-compose -f $(DOCKER_COMPOSE_FILE) -p $(APP_NAME) down --rmi local -v

.PHONY:  compose-recreate
compose-recreate: compose-destroy compose-up

.PHONY:  dev-install
dev-install: vendor
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/boumenot/gocover-cobertura@latest

.PHONY: mock-gen
mock-gen:
	cd $(CURRENT_DIR)/internal/db/repository && mockery --all --case=underscore
	cd $(CURRENT_DIR)/internal/rest/service && mockery --all --case=underscore

.PHONY: run-job
run: gen
	go run main.go task

.PHONY: automigrate
automigrate:
	go run main.go automigrate

.PHONY: migrate
migrate:
	go run main.go migrate
