project_name = service-monitor-tgbot
image_name = capcom6/$(project_name):latest

extension=
ifeq ($(OS),Windows_NT)
	extension = .exe
endif

init:
	go mod download \
		&& go install github.com/pressly/goose/v3/cmd/goose@latest \
		&& go install github.com/cosmtrek/air@latest

air:
	air

db-upgrade:
	goose up

db-upgrade-raw:
	go run ./cmd/$(project_name)/main.go db-upgrade

test:
	go test ./...

api-docs:
	swag fmt -g ./cmd/$(project_name)/main.go \
		&& swag init -g ./cmd/$(project_name)/main.go -o ./api

view-docs:
	php -S 127.0.0.1:8080 -t ./api

.PHONY: init air db-upgrade db-upgrade-raw test api-docs view-docs
