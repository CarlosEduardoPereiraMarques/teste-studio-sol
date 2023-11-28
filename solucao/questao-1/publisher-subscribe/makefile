.PHONY: test

test:
	@echo "Running tests with coverage..."
	@cd _test && go test -coverprofile=coverage.out -covermode=count -v ./...
	@go tool cover -html=_test/coverage.out -o _test/coverage.html
	@echo "Coverage report generated at: _test/coverage.html"
