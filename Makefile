include ./Makefile.vars

GO := $(shell which go)

.PHONY:
	run  \
	test \
	build

all: fmt vet build

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

run: RUN_ARGS=--help
run: fmt vet
	$(GO) run ./cmd/main.go $(RUN_ARGS)

test: generate lint
	$(GO) test ./... -cover

lint: generate
	golangci-lint run

gen: generate
generate:
	$(GO) generate ./...
