BINARY := ai
PREFIX ?= $(HOME)/.local
BINDIR := $(PREFIX)/bin

.PHONY: build install test clean

build:
	go build -o $(BINARY) .

install:
	mkdir -p $(BINDIR)
	go build -o $(BINDIR)/$(BINARY) .

test:
	go test ./...

clean:
	rm -f $(BINARY)
