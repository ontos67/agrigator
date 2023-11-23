PUBLIC_REGISTRY_HOST=docker.io
PUBLIC_REGISTRY_OWNER=ontos67
PUBLIC_REGISTRY_APP_NAME=agrigator

CI_COMMIT_REF_NAME=latest

all: deps build test

deps:
	@go mod download
	@echo "Dependencies installed successfully"

build:
	go build ./cmd/agrigator.go

test:
	go test -v ./...

lint:
	golangci-lint run ./...

image:
	docker build -t ${PUBLIC_REGISTRY_HOST}/${PUBLIC_REGISTRY_OWNER}/${PUBLIC_REGISTRY_APP_NAME}:${CI_COMMIT_REF_NAME} ./
	docker push ${PUBLIC_REGISTRY_HOST}/${PUBLIC_REGISTRY_OWNER}/${PUBLIC_REGISTRY_APP_NAME}:${CI_COMMIT_REF_NAME}