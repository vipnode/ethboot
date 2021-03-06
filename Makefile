BINARY = $(notdir $(PWD))
VERSION := $(shell git describe --tags --dirty --always 2> /dev/null || echo "dev")
SOURCES = $(wildcard *.go **/*.go)
PKG := $(shell go list | head -n1)

all: $(BINARY)

$(BINARY): # $(SOURCES)
	go build -ldflags "-X main.version=$(VERSION)" -o "$@"

deps:
	go get ./...

build: $(BINARY)

clean:
	rm $(BINARY)

run: $(BINARY)
	./$(BINARY)

debug: $(BINARY)
	./$(BINARY) --pprof :6060 -vv

test:
	go test -race ./...
	golint ./...
