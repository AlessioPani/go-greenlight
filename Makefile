# ==================================================================================== #
# CONFIGURATION
# ==================================================================================== #	
## => Edit Makefile to update configuration variables as needed
# APP CONFIGURATION VARIABLES
BINARY_NAME = greenlight
PORT = 4000
ENVIRONMENT = development

# DATABASE CONFIGURATION VARIABLES
# Only for development purpose, you should rely on your env or package like godotenv to get your dns.
DSN = "postgres://greenlight:secret_password@localhost/greenlight?sslmode=disable"
MAX_OPEN_CONNS = 25
MAX_IDLE_CONNS = 25
MAX_IDLE_TIME = 15m

# RATE LIMITER CONFIGURATION VARIABLES
LIMITER_RPS = 2
LIMITER_BURST = 4
LIMITER_ENABLED = true

# SMTP CONFIGURATION VARIABLES
SMTP_SENDER = "Greenlight <no-reply@greenlight.net>"
# MAILTRAP
#SMTP_HOST = sandbox.smtp.mailtrap.io
#SMTP_PORT = 25
#SMTP_USERNAME = <your_username>
#SMTP_PASSWORD = <your_password>
 
# MAILHOG
SMTP_HOST = localhost
SMTP_PORT =	1025

## COMMANDS LIST
# ==================================================================================== #
# HELPERS
# ==================================================================================== #	
## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #	
# build: build the application with extra flags to get the smallest executable
# -s -w : disable generation of the Go symbol table and DWARF debugging information
build:
	@echo "Building application..."
	@env go build -ldflags="-s -w" -o ./bin/api/${BINARY_NAME} cmd/api/*

# run: build and run the application
run: build
	@echo "Running application..."
	@env ./bin/api/${BINARY_NAME} -port=${PORT} -env=${ENVIRONMENT} -dsn=${DSN} -db-max-open-conns=${MAX_OPEN_CONNS} -db-max-idle-conns=${MAX_IDLE_CONNS} -db-max-idle-time=${MAX_IDLE_TIME} -limiter-rps=${LIMITER_RPS} -limiter-burst=${LIMITER_BURST} -limiter-enabled=${LIMITER_ENABLED} -smtp-host=${SMTP_HOST} -smtp-port=${SMTP_PORT} -smtp-username=${SMTP_USERNAME} -smtp-password=${SMTP_PASSWORD} -smtp-sender=${SMTP_SENDER}

## start: starts the application
start: run

## stop: stops the running application
# Windows users: use @taskkill /IM ${BINARY_NAME} /F instead
stop:
	@echo "Stopping..."
	@-pkill -SIGTERM -f "${BINARY_NAME}"

## restart: stop and start the application
restart: stop start

## migrate-up: executes db migrations
migrate-up:
	@env migrate -path=./migrations -database=${DSN} up

## migrate-down: revert db migrations
migrate-down:
	@env migrate -path=./migrations -database=${DSN} down

## migration: create a new set of migration files (make migration name=<migration_name>)
migration:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## clean: runs go clean and deletes the executable
clean:
	@echo "Cleaning..."
	@go clean -testcache
	@-rm ./bin/api/${BINARY_NAME}

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #	
## test: executes tests in verbose mode
test:
	@env go test -vet=off -p 1 ./...

## coverage: executes tests and generate coverage profile
coverage:
	@env go test ./... -coverprofile=./coverage.out  -coverpkg=./... && go tool cover -html=./coverage.out

## tidy: format all .go files and tidy module dependencies
tidy:
	@echo 'Formatting .go files...'
	@env go fmt ./...
	@echo 'Tidying module dependencies...'
	@env go mod tidy
	@echo 'Verifying and vendoring module dependencies...'
	go mod verify
	go mod vendor

## audit: run quality control checks
audit:
	@echo 'Checking module dependencies'
	@env go mod tidy -diff
	@env go mod verify
	@echo 'Vetting code...'
	@env go vet ./...
	@echo 'Running tests...'
	@env go test -vet=off -v ./...

