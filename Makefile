SHELL=/bin/bash

run:
	source .env && swag init -g cmd/main.go --parseDependency && air --build.cmd "go build -o bin/api cmd/main.go" --build.bin "./bin/api"

test:
	go test -v ./internal/core/services/...