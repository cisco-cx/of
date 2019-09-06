.PHONY: format
format:
	gofmt -w .

.PHONY: test
test:
	go test ./...
