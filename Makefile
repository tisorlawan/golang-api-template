# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

DOCKER_NAME=$(shell basename $(shell pwd)):latest
BINARY=main

all: test build
build:
	$(GOBUILD) -o $(BINARY)
test:
	$(GOTEST) ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY)
run: clean build
	./$(BINARY)
deps:
	$(GOGET)

# Single binary build
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY)

# Docker related task
docker-build: clean build-linux
	docker build -t $(DOCKER_NAME) .
