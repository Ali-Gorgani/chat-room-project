run-auth:
	cd services/auth-service && go run cmd/main.go

run-user:
	cd services/user-management && go run cmd/main.go

run-chat:
	cd services/chat-service && go run cmd/main.go

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_auth build_user build_chat
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_auth: builds the auth binary as a linux executable
build_auth:
	@echo "Building auth binary..."
	cd ./services/auth-service && env GOOS=linux CGO_ENABLED=0 go build -o bin/auth-service cmd/main.go && chmod +x bin/auth-service
	@echo "Done!"

## build_user: builds the user binary as a linux executable
build_user:
	@echo "Building user binary..."
	cd ./services/user-management && env GOOS=linux CGO_ENABLED=0 go build -o bin/user-service cmd/main.go && chmod +x bin/user-service
	@echo "Done!"

## build_chat: builds the chat binary as a linux executable
build_chat:
	@echo "Building chat binary..."
	cd ./services/chat-service && env GOOS=linux CGO_ENABLED=0 go build -o bin/chat-service cmd/main.go && chmod +x bin/chat-service
	@echo "Done!"

.PHONY: run-auth run-user run-chat up up_build down build_auth build_user build_chat