.PHONY: wire, up, down, swag
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