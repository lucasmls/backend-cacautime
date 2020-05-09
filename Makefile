
include .env
export $(shell sed 's/=.*//' .env)

GOPATH=$(shell go env GOPATH)

deps:
	@ echo
	@ echo "Starting downloading dependencies..."
	@ echo
	@ go get -u ./...

server:
	@ echo
	@ echo "Starting the server..."
	@ echo
	@ go run ./cmd/server