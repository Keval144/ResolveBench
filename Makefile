BINARY ?= resolvebench
GO ?= go

.PHONY: build run test clean lint fmt install

build:
	$(GO) build -o $(BINARY) .

run:
	$(GO) run .

test:
	$(GO) test -v -race ./...

clean:
	$(GO) clean
	rm -f $(BINARY)

lint:
	$(GO) vet ./...

fmt:
	$(GO) fmt ./...

install: build
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	cp $(BINARY) $(DESTDIR)$(PREFIX)/bin/
