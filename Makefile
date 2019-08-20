SHELL := /bin/bash

ROOT := $$(git rev-parse --show-toplevel)

build:
	go build -o $(ROOT)/bin/pub $(ROOT)/cmd/pub/main.go
	go build -o $(ROOT)/bin/sub $(ROOT)/cmd/sub/main.go
