SHELL:=/bin/bash

.PHONY: golint-ci
golint-ci:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v2.0.1 golangci-lint run -v --timeout=5m


.PHONY: test
test:
	go test -v -race ./...

.PHONY: test-coverage
test-coverage:
	go test -test.v -race -coverprofile=coverage.txt -covermode=atomic ./...
