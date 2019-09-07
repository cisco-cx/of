.PHONY: refine
refine:  ## Run all formatters and update dependency list.
	gofmt -w .
	go mod tidy

.PHONY: test
test:
	go test ./...
