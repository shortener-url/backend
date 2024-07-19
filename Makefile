SHELL=/bin/bash

run:
	source .env && swag init --parseDependency && air --build.cmd "go build -o bin/api main.go" --build.bin "./bin/api"