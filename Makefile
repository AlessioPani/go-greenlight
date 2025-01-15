## APP CONFIGURATION VARIABLES
BINARY_NAME = greenlight
PORT = 4000
ENVIRONMENT = development
## DATABASE CONFIGURATION VARIABLES
# Only for development purpose, you should rely on your env or package like godotenv to get your dns.
DSN = "postgres://greenlight:secret_password@localhost/greenlight?sslmode=disable"
MAX_OPEN_CONNS = 25
MAX_IDLE_CONNS = 25
MAX_IDLE_TIME = 15m

## COMMANDS LIST
# build: build the application with extra flags to get the smallest executable
# -s -w : disable generation of the Go symbol table and DWARF debugging information
build:
	@echo "Building application..."
	@env go build -ldflags="-s -w" -o ${BINARY_NAME} cmd/api/*

# run: build and run the application
run: build
	@echo "Running application..."
	@env ./${BINARY_NAME} -port=${PORT} -env=${ENVIRONMENT} -dsn=${DSN} -db-max-open-conns=${MAX_OPEN_CONNS} -db-max-idle-conns=${MAX_IDLE_CONNS} -db-max-idle-time=${MAX_IDLE_TIME}

# start: alias to run
start: run

# stop: stops the running application
# Windows users: use @taskkill /IM ${BINARY_NAME} /F instead
stop:
	@echo "Stopping..."
	@-pkill -SIGTERM -f "${BINARY_NAME}"

# restart: stop and start the application
restart: stop start

# test: executes tests in verbose mode
test:
	@env go test -v ./...

# coverage: executes tests and generate coverage profile
coverage:
	@env go test -coverprofile=./coverage.out  ./... && go tool cover -html=./coverage.out

# clean: runs go clean and deletes the executable
clean:
	@echo "Cleaning..."
	@go clean
	@rm ${BINARY_NAME}
