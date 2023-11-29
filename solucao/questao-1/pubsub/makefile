.PHONY: test

test:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out -covermode=count -v ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at: coverage.html"
