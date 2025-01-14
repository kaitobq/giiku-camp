.PHONY: wire, up, down, swag, test, test-coverage, mock
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
mock:
	mockgen -source=internal/usecase/usecase.go -destination=internal/mock/usecase/usecase_mock.go -package=mock_usecase && \
	mockgen -source=internal/domain/repository/user_repo.go -destination=internal/mock/domain/repository/user_repo_mock.go -package=mock_repository