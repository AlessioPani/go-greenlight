## CONFIGURATION VARIABLES
BINARY_NAME = greenlight
PORT = 4000
ENVIRONMENT = development
# Only for development purpose, you should rely on your env or package like godotenv.
DSN = "postgres://greenlight:secret_password@localhost/greenlight?sslmode=disable"

## COMMANDS LIST
# build: build the application with extra flags to get the smallest executable
# -s -w : disable generation of the Go symbol table and DWARF debugging information
build:
	@echo "Building application..."
	@env go build -ldflags="-s -w" -o ${BINARY_NAME} cmd/api/*
	@echo "...done!"

# run: build and run the application
run: build
	@echo "Running application..."
	@env ./${BINARY_NAME} -port=${PORT} -env=${ENVIRONMENT} -dsn=${DSN}

# start: alias to run
start: run

# stop: stops the running application
# Windows users: use @taskkill /IM ${BINARY_NAME} /F instead
stop:
	@echo "Stopping..."
	@-pkill -SIGTERM -f "${BINARY_NAME}"
	@echo "...done!"

# restart: stop and start the application
restart: stop start

# clean: runs go clean and deletes the executable
clean:
	@echo "Cleaning..."
	@go clean
	@rm ${BINARY_NAME}
	@echo "...done!"
