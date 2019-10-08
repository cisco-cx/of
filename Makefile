.DEFAULT_GOAL := help
SHELL := bash

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
all: | refine test build  ## Run all generally applicable targets.

.PHONY: build
build:  ## Build the program for Linux.
	@echo "==> Building the program for Linux."
	CGO_ENABLED=0 GOOS=linux go build -v -a -ldflags '$(LDFLAGS)' -mod=vendor .

.PHONY: clean
clean:  ## Clean temporary files.
	@echo "==> Cleaning temporary files."
	rm -f ./*.html
	rm -f ./*.pdf
	rm -f ./*.pprof
	rm -f ./cp.out
	rm -f ./of

.PHONY: docker
docker:  ## Build a docker image for local dev.
	docker build . -t of:local

.PHONY: demo-docker
demo-docker:  ## Run the demo using Docker image
	# TODO: Use docker-compose and add Alertmanager.
	docker run --rm --net=host -v /keybase/path -it of:local demo

.PHONY: refine
refine:  ## Run all formatters and static analysis.
	@echo "==> Running all formatters and static analysis."
	gofmt -w .

.PHONY: report
report:  ## Generate all reports.
	@echo "==> Generating profiler reports."
	for mode in cpu mem mutex block; do \
	  if [ -e $$mode.pprof ]; then go tool pprof --pdf $(PROGRAM) $$mode.pprof > $$mode.pprof.pdf; fi; done
	@if [ -d ~/x/tmp ] && compgen -G "*.pdf"; then cp -v *.pdf ~/x/tmp; fi
	@echo "==> Generating coverage reports."
	go tool cover -html=cp.out -o=coverage.html

.PHONY: test
test:  ## Run all tests and generate all reports.
	@echo "==> Running all tests."
	go test ./... -coverprofile=cp.out -mod=vendor
	go tool cover -func=cp.out
	@$(MAKE) vet
	@$(MAKE) report

.PHONY: tidy
tidy:  ## Run go mod tidy (depends on access to github.com)
	@echo "==> Tidying dependency list."
	go mod tidy -v

.PHONY: vet
vet:  ## Run Go vet.
	@echo "==> Running Go vet."
	go vet ./...

.PHONY: vendor
vendor:  ## Re-vendor dependencies. (depends on access to github.com)
	@echo "==> Re-vendoring dependencies."
	go mod vendor -v
	go mod verify

.PHONY: help
help:  ## Print list of Makefile targets.
	@# Taken from https://github.com/spf13/hugo/blob/master/Makefile
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	  cut -d ":" -f1- | \
	  awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
