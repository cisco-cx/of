.PHONY: all
all: | refine test

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
