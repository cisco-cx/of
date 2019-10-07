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

.PHONY: clean
clean:  ## Clean temporary files.
	@echo "==> Cleaning temporary files."
	rm -f ./*.html
	rm -f ./*.pdf
	rm -f ./*.pprof
	rm -f ./cp.out

.PHONY: refine
refine:  ## Run all formatters and static analysis. Tidy dependency list.
	@echo "==> Running all formatters and static analysis."
	gofmt -w .
	@echo "==> Tidying dependency list."
	go mod tidy

.PHONY: report
report:  ## Generate all reports.
	@echo "==> Generating profiler reports."
	for mode in cpu mem mutex block; do \
	  if [ -e $$mode.pprof ]; then go tool pprof --pdf $(PROGRAM) $$mode.pprof > $$mode.pprof.pdf; fi; done
	@if [ -d ~/x/tmp ]; then cp -v *pdf ~/x/tmp; fi
	@echo "==> Generating coverage reports."
	go tool cover -html=cp.out -o=coverage.html

.PHONY: test
test:  ## Run all tests and generate all reports.
	@echo "==> Running all tests."
	go test ./... -coverprofile=cp.out
	@$(MAKE) reports
