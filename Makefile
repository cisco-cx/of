# based in part on: https://povilasv.me/exposing-go-modules-to-prometheus/
PROGRAM := of
PACKAGE := github.com/cisco-cx/$(PROGRAM)
INFO_PACKAGE := $(PACKAGE)/info
LICENSE := Apache-2.0
URL     := https://$(PACKAGE)
DATE := $(shell date +%FT%T%z)
USER := $(shell whoami)
GIT_HASH := $(shell git --no-pager describe --tags --always)
BRANCH := $(shell git branch | grep '*' | cut -d ' ' -f2)

LDFLAGS := -s
LDFLAGS += -X "$(INFO_PACKAGE).Program=$(PROGRAM)"
LDFLAGS += -X "$(INFO_PACKAGE).License=$(LICENSE)"
LDFLAGS += -X "$(INFO_PACKAGE).URL=$(URL)"
LDFLAGS += -X "$(INFO_PACKAGE).BuildUser=$(USER)"
LDFLAGS += -X "$(INFO_PACKAGE).BuildDate=$(DATE)"
LDFLAGS += -X "$(INFO_PACKAGE).Version=$(GIT_HASH)"
LDFLAGS += -X "$(INFO_PACKAGE).Revision=$(GIT_HASH)"
LDFLAGS += -X "$(INFO_PACKAGE).Branch=$(BRANCH)"

.PHONY: all
all: | refine test build

.PHONY: build
build:  ## Build the program for Linux.
	@echo "==> Building the program for Linux."
	CGO_ENABLED=0 GOOS=linux go build -v -a -ldflags '$(LDFLAGS)' .

.PHONY: refine
refine:  ## Run all formatters and static analysis. Tidy dependency list.
	@echo "==> Running all formatters and static analysis."
	gofmt -w .
	@echo "==> Tidying dependency list."
	go mod tidy

.PHONY: test
test:  ## Run all tests. Generate 'coverage.html'.
	@echo "==> Running all tests."
	go test ./... -coverprofile=cp.out
	@echo "==> Generating 'coverage.html'."
	go tool cover -html=cp.out -o=coverage.html
