# service name
APP_NAME := skeleton

# The binary to build (just the basename).
BINARY ?= skeleton

# This repo's root import path (under GOPATH).
PKG := github.com/screwyprof/skeleton

## DO NOT EDIT LINES BELOW
GO_FILES = $(shell find . -name "*.go" | grep -v vendor | uniq)

BUILD_DATE = $(shell date -u '+%Y.%m.%d')
GOLANG_VERSION ?= $(shell go version | cut -d" " -f3 | sed 's/go//')

GIT_REV    ?= $(shell git rev-parse --short HEAD)
GIT_TAG    ?= $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_BRANCH ?= $(shell git branch|grep '*'| cut -f2 -d' ')
GIT_DIRTY  ?= $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")
GIT_LOG    ?= $(shell git log --decorate --oneline -n1| sed -e "s/'/ /g" -e "s/\"/ /g" -e 's/`/ /g')

ifdef SOURCE_VERSION
	GIT_REV = $(SOURCE_VERSION)
	BINARY_VERSION = $(BUILD_DATE)-$(SOURCE_VERSION)
endif

BINARY_VERSION  ?= ${BUILD_DATE}-${GIT_BRANCH}-${GIT_REV}

# Only set Version if building a tag or VERSION is set
ifneq ($(BINARY_VERSION),)
	LDFLAGS += -X $(PKG)/internal/pkg/app/version.AppVersion=$(VERSION)
endif

LDFLAGS = -X $(PKG)/internal/pkg/app/version.AppName=$(APP_NAME)
LDFLAGS += -X $(PKG)/internal/pkg/app/version.AppVersion=$(BINARY_VERSION)
LDFLAGS += -X $(PKG)/internal/pkg/app/version.GitRev=$(GIT_REV)
LDFLAGS += -X $(PKG)/internal/pkg/app/version.GoVersion=$(GOLANG_VERSION)
LDFLAGS += -X $(PKG)/internal/pkg/app/version.BuildDate=$(BUILD_DATE)
LDFLAGS += -X '$(PKG)/internal/pkg/app/version.GitLog=$(GIT_LOG)'
#LDFLAGS += -X main.GitTreeState=${GIT_DIRTY}

OK_COLOR=\033[32;01m
NO_COLOR=\033[0m
MAKE_COLOR=\033[33;01m%-20s\033[0m
        
all:  install-tools lint build test ## build, lint and test

ci-all: install-tools lint build test-ci ## run build, lint and test pipeline

deps: ## sync go mod deps
	@echo "$(OK_COLOR)--> Download go.mod dependencies$(NO_COLOR)"
	go mod download

install-tools: ## install dev tools, linters, code generaters, etc
	@echo "$(OK_COLOR)--> Installing tools from tools/tools.go$(NO_COLOR)"
	@export GOBIN=$$PWD/tools/bin; export PATH=$$GOBIN:$$PATH; cat tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

build: ## build application
	@echo "$(OK_COLOR)--> Building application$(NO_COLOR)"
	go build -ldflags "$(LDFLAGS)" -o $(PWD)/bin/$(BINARY) $(PWD)/cmd/skeleton.go

build-docker: ## build application staically for docker
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -ldflags "-w $(LDFLAGS)" -a -o ./bin/$(BINARY)  $(PWD)/cmd/skeleton.go

build-ci: ## build application with race detector
	@echo "$(OK_COLOR)--> Building application$(NO_COLOR)"
	go build -race -ldflags "$(LDFLAGS)" -o $(PWD)/bin/$(BINARY) $(PWD)/cmd/skeleton.go

run: # run application locally with the given .env file
	@echo "$(OK_COLOR)--> Running application$(NO_COLOR)"
	@(sh -ac 'source .env && go run cmd/skeleton.go')

lint: ## run linters
	@echo "$(OK_COLOR)--> Running linters$(NO_COLOR)"
	tools/bin/golangci-lint run

mock-gen: ## generate mocks
	@echo "$(OK_COLOR)--> Generating mocks$(NO_COLOR)"
	tools/bin/mockgen -source=pkg/cert/usecase/viewcert/view_certificate.go -package=mock -destination=pkg/cert/mock/cert_reporter_mock.go
	tools/bin/mockgen -source=pkg/cert/usecase/issuecert/issue_certificate.go -package=mock -destination=pkg/cert/mock/cert_storage_mock.go
	tools/bin/mockgen -source=vendor/github.com/screwyprof/golibs/queryer/queryer.go -package mock -destination=internal/pkg/delivery/rest/mock/queryer_mock.go
	tools/bin/mockgen -source=vendor/github.com/screwyprof/golibs/cmdhandler/command_handler.go -package mock -destination=internal/pkg/delivery/rest/mock/command_handler_mock.go
             
test: test-unit test-integration test-e2e # run all tests

test-docker: ## run all tests in with docker-compose
	docker-compose up --build --force-recreate --no-deps -d db
	docker-compose run migrate
	@(sh -ac 'source .env && make test')
	docker-compose down --remove-orphans --volumes
	
test-unit: ## run unit tests
	@echo "$(OK_COLOR)--> Running unit tests$(NO_COLOR)"
	go test --race --count=1 ./...

test-integration: ## run integration tests
	@echo "$(OK_COLOR)--> Running integration tests$(NO_COLOR)"
	go test --tags "integration" --race --count=1 ./tests/integration/...

test-e2e: ## run e2e tests
	@echo "$(OK_COLOR)--> Running E2E tests$(NO_COLOR)"
	go test --tags "acceptance" --race --count=1 ./tests/e2e/...

test-ci: ## runing all tests with coverage
	@echo "$(OK_COLOR)--> Generating code coverage$(NO_COLOR)"
	tools/generate-fake-tests.sh
	tools/coverage.sh

fmt: ## format go files
	@echo "$(OK_COLOR)--> Formatting go files$(NO_COLOR)"
	@tools/bin/gofumpt -s -l -w $(GO_FILES)
	@tools/bin/gci -w -l $(PKG) $(GO_FILES)

clean: ## clean up
	@echo "$(OK_COLOR)--> Clean up$(NO_COLOR)"
	rm -rf $(PWD)/tools/bin
	rm -rf $(PWD)/bin/$(BINARY)

version: ## show build info
	@echo "$(OK_COLOR)--> Build info$(NO_COLOR)"
	@echo "Version:           ${BINARY_VERSION}"
	@echo "Date:              ${BUILD_DATE}"
	@echo "Git Tag:           ${GIT_TAG}"
	@echo "Git Rev:           ${GIT_REV}"
	@echo "Git Tree State:    ${GIT_DIRTY}"

help: ## show this help
	@egrep '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(MAKE_COLOR) %s\n", $$1, $$2}'

# To avoid unintended conflicts with file names, always add to .PHONY
# unless there is a reason not to.
# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
.PHONY: all ci-all deps install-tools build build-ci run lint mock-gen test test-local test-unit test-integration test-e2e test-ci fmt clean version help