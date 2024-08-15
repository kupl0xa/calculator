#############
# VARIABLES #
#############

BINARY_NAME=calculator
DOCKERFILE_PATH=./Dockerfile
DOCKER_IMAGE_NAME=calculator-image
DOCKER_IMAGE_TAG=latest

#############
# COMMANDS  #
#############

download:
	@echo Download go.mod dependencies
	@go mod download

build:
	@echo "Building the project"
	@go build -o $(BINARY_NAME) cmd/calculator/main.go

run: build
	@echo "Starting the server"
	@./$(BINARY_NAME)

test:
	@echo "Running tests"
	@go test ./...

test-cov:
	@echo "Tests coverage"
	@go test -cover ./...
	
cov-report:
	@echo "Tests coverage report"
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@rm -f coverage.out

lint:
	@go mod tidy
	@golangci-lint run

docker-build:
	@echo "building Docker image"
	@docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) -f $(DOCKERFILE_PATH) .

docker-run: docker-build
	@echo "Starting Docker container"
	@docker run -p 8080:8080 $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

clean:
	@echo "Cleaning up"
	@go clean
	@rm -f $(BINARY_NAME)
