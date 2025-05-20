# ==============================================
# Makefile variables
# ==============================================
APP_BINARY=app
APP_NAME=microsoft-apps-exporter

CACHE_DIR=./cache
TESTS_CACHE_DIR=${CACHE_DIR}/tests

MIGRATIONS_DIR=./migrations
MIGRATION_TABLE=microsoft_apps_exporter_migrations

TEST_HELM_VALUES_PATH=./helm/${APP_NAME}/values-minikube.yaml

# ==============================================
# Environment Management
# ==============================================

## load_development_env: Copies .env.development to .env and exports variables
load_development_env:
	@echo "📋 Loading development environment variables..."
	@cp .env.development .env
	@echo "✅ Testing environment loaded (secrets files override environment variables)." 

## load_testing_env: Copies .env.testing to .env and exports variables
load_testing_env:
	@echo "📋 Loading testing environment variables..."
	@cp .env.testing .env
	@echo "✅ Testing environment loaded (secrets files override environment variables)." 

# ==============================================
# Docker Management
# ==============================================

## up_build: Stops containers, builds everything, and starts production containers (No migrations)
up_build: down build
	@echo "✅ Environment is ready!"

## build: Builds and starts all containers
build:
	@echo "🛠️ Building and starting containers..."
	docker-compose up --build -d
	@echo "✅ Containers built and started!"

## up: Starts only production containers without rebuilding images
up:
	@echo "▶️ Starting containers..."
	docker-compose up -d
	@echo "✅ Containers started!"

## down: Stops all containers
down:
	@echo "🛑 Stopping all containers..."
	docker-compose down
	@echo "✅ All containers stopped!"

## restart: Restarts the application container
restart:
	@echo "🔄 Restarting application container..."
	docker-compose restart $(APP_NAME)
	@echo "✅ Application container restarted!"

## logs: Displays logs from the application container
logs:
	@echo "📜 Displaying logs..."
	docker-compose logs -f $(APP_NAME)

# ==============================================
# Build & Run Commands
# ==============================================

## dry_run: Runs the application locally
dry_run:
	@echo "🚀 Running application locally..."
	go run ./cmd/main.go

## binary_build: Builds the application binary for Linux
binary_build:
	@echo "🛠️ Building app binary..."
	go mod tidy
	env GOOS=linux CGO_ENABLED=0 go build -o $(APP_BINARY) ./cmd/main.go
	@echo "✅ App binary built successfully!"

## binary_stop: Stops the running application process
binary_stop:
	@echo "🛑 Stopping application..."
	@-pkill -SIGTERM -f "$(APP_BINARY)" || echo "No running process found."
	@echo "✅ Application stopped!"

# ==============================================
# Database & Migrations (For Local Dev Only)
# ==============================================

## create_migration: Creates a new migration file (Usage: `make create_migration name=your_migration`)
create_migration:
	@echo "📄 Creating migration: $(name)..."
	goose create "$(name)" sql --dir=$(MIGRATIONS_DIR) --table=$(MIGRATION_TABLE)

## down_migration: Rolls back the last migration
down_migration:
	@echo "⏪ Rolling back the last migration..."
	goose postgres "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" \
		down --dir=$(MIGRATIONS_DIR) --table=$(MIGRATION_TABLE)
	@echo "✅ Migration rolled back!"

## migrate: Runs all pending migrations (For local development only)
migrate:
	@echo "📊 Applying database migrations (local only)..."
	goose postgres "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" \
		up --dir=$(MIGRATIONS_DIR) --table=$(MIGRATION_TABLE)
	@echo "✅ Migrations applied successfully!"

# ==============================================
# Testing
# ==============================================

## unit_testing: Runs unit tests with coverage reporting
unit_testing:
	@echo "🧪 Running unit tests..."
	go test ./tests/unit/... -v -count=1 -tags="testing unit" -coverprofile=${TESTS_CACHE_DIR}/coverage.unit.out -coverpkg=./...
	@echo "✅ Unit tests completed!"

## integration_testing: Runs integration tests with coverage reporting
integration_testing: down
	@echo "🧪 Running integration tests..."
	docker-compose up postgres pgadmin -d
	go test ./tests/integration/... -v -count=1 -tags="testing integration" -coverprofile=${TESTS_CACHE_DIR}/coverage.integration.out -coverpkg=./...
	@echo "✅ Integrations tests completed!"

## e2e_api_testing: Runs e2e tests of api package with coverage reporting
e2e_api_testing: down
	@echo "🧪 Running e2e tests..."
	go test ./tests/e2e/... -v -count=1 -tags="testing e2e" -coverprofile=${TESTS_CACHE_DIR}/coverage.e2e.out -coverpkg=./...
	@echo "✅ e2e tests completed!"

## merge_coverage: Merge all test coverage reports and generate visual reports
## Generates both HTML and SVG treemap coverage visualizations
merge_coverage:
	@echo "🔀 Merging coverage reports..."
	gocovmerge ${TESTS_CACHE_DIR}/coverage.unit.out \
		${TESTS_CACHE_DIR}/coverage.integration.out \
		${TESTS_CACHE_DIR}/coverage.e2e.out > ${TESTS_CACHE_DIR}/coverage.merged.out

	@echo "📊 Generating final coverage reports..."
	go tool cover -html=${TESTS_CACHE_DIR}/coverage.merged.out -o ${TESTS_CACHE_DIR}/coverage.html
	go-cover-treemap -coverprofile ${TESTS_CACHE_DIR}/coverage.merged.out > ${TESTS_CACHE_DIR}/coverage.svg
	@echo "✅ Coverage report is available in ${TESTS_CACHE_DIR}/coverage.svg"

# ==============================================
# Minikube Management
# ==============================================

## minikube_start: Starts Minikube with specified resources
minikube_start:
	@echo "🚀 Starting Minikube..."
	minikube start --cpus=2 --memory=4g --driver=docker
	minikube addons enable ingress
	@echo "✅ Minikube started!"

## minikube_stop: Stops Minikube
minikube_stop:
	@echo "🛑 Stopping Minikube..."
	minikube stop
	@echo "✅ Minikube stopped!"

## minikube_clean: Deletes Minikube cluster and starts fresh
minikube_clean:
	@echo "🧹 Cleaning Minikube..."
	minikube delete
	minikube start --cpus=2 --memory=4g --driver=docker
	minikube addons enable ingress
	@echo "✅ Minikube cleaned and restarted!"

## minikube_monkeybot_clean: Cleans up monkeybot namespace
minikube_monkeybot_clean:
	@echo "🧹 Cleaning monkeybot namespace..."
	kubectl delete namespace monkeybot || true
	kubectl create namespace monkeybot
	@echo "✅ monkeybot namespace cleaned!"

## minikube_image_build: Builds and loads app image to Minikube
minikube_image_build:
	@echo "🛠️ Building and loading image..."
	docker build -t ${APP_NAME}:latest .
	minikube ssh -- docker rmi -f ${APP_NAME}:latest || true
	minikube image load ${APP_NAME}:latest
	@echo "✅ Image built and loaded!"

## minikube_deploy_all: Deploys both db-monkeybot and app
minikube_deploy_all: minikube_monkeybot_clean minikube_image_build
	@echo "🚀 Deploying all services..."
	helm install db-monkeybot ./helm/db-monkeybot -n monkeybot
	helm install ${APP_NAME} ./helm/${APP_NAME} -n monkeybot -f ${TEST_HELM_VALUES_PATH}
	kubectl -n monkeybot port-forward svc/${APP_NAME} 8080:8080
	@echo "✅ All services deployed!"

## minikube_deploy_app: Updates app (code or chart changes)
minikube_deploy_app:
	@echo "🚀 Deploying ${APP_NAME}..."
	helm upgrade --install ${APP_NAME} ./helm/${APP_NAME} -n monkeybot -f ${TEST_HELM_VALUES_PATH}
	kubectl -n monkeybot port-forward svc/${APP_NAME} 8080:8080
	@echo "✅ ${APP_NAME} deployed!"

## minikube_update_app_code: Updates and deploys app source code
minikube_update_app_code: minikube_image_build
	@echo "🚀 Updating ${APP_NAME} source code..."
	kubectl -n monkeybot rollout restart deployment ${APP_NAME}
	@echo "✅ Source code updated and deployed!"

## minikube_logs_app: Shows logs for app
minikube_logs_app:
	kubectl -n monkeybot logs -l app=${APP_NAME}

## minikube_get_all: Lists all Kubernetes resources in the 'monkeybot' namespace
minikube_get_all:
	kubectl -n monkeybot get all

## minikube_delete_app: Deletes the app pod to trigger a restart
minikube_delete_app:
	kubectl -n monkeybot delete pod -l app=${APP_NAME}
