.PHONY: test

test:
	@echo "Running tests..."
	@cd pubsub && go test -coverprofile=coverage.out -covermode=count -v ./...
	@cd pubsub && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at: pubsub/coverage.html"

run:
	@echo "Running application..."
	@go run main.go
