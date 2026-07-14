.PHONY: run swagger dev-run mock test
# Run the application
run:
	go run cmd/api/main.go

# Generate/update Swagger API docs
swagger:
	swag init -g cmd/api/main.go

# Regenerate docs then run (useful during development)
dev-run: swagger run

# Generate mocks for interfaces (mockery, driven by //go:generate comments)
mock:
	go generate ./internal/service/...
	go generate ./internal/repository/...

# Coverage config
COVERAGE_EXCLUDE=mocks|main.go|test
COVERAGE_THRESHOLD=80

# Run all tests with coverage report + enforce threshold
test:
	go test ./... -coverprofile=coverage.tmp -covermode=atomic -coverpkg=./... -p 1
	grep -vE "$(COVERAGE_EXCLUDE)" coverage.tmp > coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@total=$$(go tool cover -func=coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "❌ Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
		exit 1; \
	else \
		echo "✅ Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
	fi