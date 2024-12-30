.PHONY: wire, up, down
wire:
	cd internal/app && go run -mod=mod github.com/google/wire/cmd/wire
up:
	docker-compose up --watch
down:
	docker-compose down
