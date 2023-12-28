project_name = service-monitor-tgbot
image_name = capcom6/$(project_name):latest

extension=
ifeq ($(OS),Windows_NT)
	extension = .exe
endif

run:
	go run ./cmd/$(project_name)/main.go

monitor:
	CONFIG_PATH=./configs/monitor.yml go run ./cmd/monitor/main.go

bot:
	CONFIG_PATH=./configs/bot.yml go run ./cmd/bot/main.go

init-dev:
	go mod download \
		&& go install github.com/cosmtrek/air@latest

init: init-dev
	go mod download \
		&& go install github.com/cosmtrek/air@latest

air-monitor:
	CONFIG_PATH=./configs/monitor.yml air -c monitor.air.toml

air-bot:
	CONFIG_PATH=./configs/bot.yml air -c bot.air.toml

test:
	go test ./...

api-docs:
	swag fmt -g ./cmd/$(project_name)/main.go \
		&& swag init -g ./cmd/$(project_name)/main.go -o ./api

view-docs:
	php -S 127.0.0.1:8080 -t ./api

docker-build:
	docker build -f build/package/Dockerfile -t capcom6/service-monior-core --build-arg APP=monitor .
	docker build -f build/package/Dockerfile -t capcom6/service-monior-bot --build-arg APP=bot .

docker:
	docker-compose -f deployments/docker-compose/docker-compose.yml up --build

docker-dev:
	docker-compose -f deployments/docker-compose/docker-compose.dev.yml up --build

.PHONY: run monitor bot init-dev init air-monitor air-bot test api-docs view-docs docker-build docker
