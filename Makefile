MAIN_BINARY=main
DOCKER_COMPOSE=docker-compose

## up: starts only production containers in the background without forcing build
up:
	@echo "Starting Docker images for production..."
	$(DOCKER_COMPOSE) up -d
	@echo "Production containers started!"

## up_dev: starts development containers in the background (includes Postgres & PGAdmin)
up_dev:
	@echo "Starting Docker images for development..."
	$(DOCKER_COMPOSE) --profile dev up -d
	@echo "Development containers started!"

## up_build: stops containers, builds everything, and starts production containers
up_build: build_main
	@echo "Stopping running containers..."
	$(DOCKER_COMPOSE) down
	@echo "Building and starting Docker images for production..."
	$(DOCKER_COMPOSE) up --build -d
	@echo "Production containers built and started!"

## up_build_dev: stops containers, builds everything, and starts development containers
up_build_dev: build_main
	@echo "Stopping running containers..."
	$(DOCKER_COMPOSE) down
	@echo "Building and starting Docker images for development..."
	$(DOCKER_COMPOSE) --profile dev up --build -d
	@echo "Development containers built and started!"

## down: stops all containers
down:
	@echo "Stopping all containers..."
	$(DOCKER_COMPOSE) --profile dev down || $(DOCKER_COMPOSE) down
	@echo "All containers stopped!"

## build_main: builds the main binary as a Linux executable
build_main: cleanup
	@echo "Building main binary..."
	go mod tidy
	env GOOS=linux CGO_ENABLED=0 go build -o ./cache/${MAIN_BINARY} ./cmd
	@echo "Main binary built successfully!"

## cleanup: removes cached files
cleanup:
	@echo "Cleaning up cache directory..."
	sudo rm -rf ./cache/*
	@echo "Cache directory cleaned!"

## run: runs the main application
run:
	@echo "Starting main application..."
	./entrypoint.sh

## stop: stops the main application gracefully
stop:
	@echo "Stopping main application..."
	@-pkill -SIGTERM -f "./cache/${MAIN_BINARY}"
	@echo "Main application stopped!"

## test_unit: runs unit tests with coverage and verbose output
test_unit:
	@echo "Running unit tests..."
	go test -coverprofile=coverage.out ./tests/... -v -tags "testing unit"
	@echo "Unit tests completed!"
	@echo "Generating coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
