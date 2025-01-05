.PHONY: wire, up, down, swag, test, test-coverage
wire:
	cd internal/app && go run -mod=mod github.com/google/wire/cmd/wire
up:
	docker-compose up --watch
down:
	docker-compose down
swag:
	swag init \
	--parseDependency \
	--parseInternal \
	--parseDepth 5 \
	--generalInfo ./cmd/server/main.go \
	--dir .
test:
	go test -v ./...
cover:
	go test -v -coverprofile=coverage/coverage.out ./... && \
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html && \
	go tool cover -func=coverage/coverage.out > coverage/coverage.txt
	open coverage/coverage.html