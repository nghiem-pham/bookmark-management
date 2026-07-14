run:
	go run ./cmd/api/main.go

test:
	go test ./... -v -count=1

mock:
	go generate ./internal/service/...