# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
BINARY_NAME=marumaru
CLUSTER_NAME=marumaru
# DB_PASSWORD = ${MONGODB_PASSWORD}

all: deps build test clean

build:
		GO111MODULE=on CGO_ENABLED=0 $(GOBUILD) -o $(BINARY_NAME) cmd/marumaru-backend/marumaru-backend.go
test:
		$(GOTEST) -v ./server
clean:
		$(GOCLEAN)
		rm -f ./$(BINARY_NAME)
run:
		$(GORUN) cmd/marumaru-backend/marumaru-backend.go
deps:
		$(GOMOD) download
