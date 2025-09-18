SHELL := /bin/bash
APP := bff
PKG := ./cmd/bff

# Auto-load .env if present
ifneq (,$(wildcard .env))
include .env
export
endif

.PHONY: dev build docker tidy

dev:
	@echo "Starting with PORT=$(PORT)  UPSTREAM_SIF=$(UPSTREAM_SIF)  UPSTREAM_PAYMENTS=$(UPSTREAM_PAYMENTS)"
	go run $(PKG)

build:
	CGO_ENABLED=0 go build -o bin/$(APP) $(PKG)

docker:
	docker build -t bff-api:local .

tidy:
	go mod tidy
