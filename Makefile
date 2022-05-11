SHELL = /bin/bash

test:
	go test ./...

run:
	source ".env.local"
	go run cmd/subscriber/main.go
