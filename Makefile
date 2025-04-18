PREFIX ?= /usr/local

build:
	[ -d bin ] || mkdir bin
	go build -o bin .

install:
	install bin/* $(PREFIX)/bin

test:
	go test

clean:
	[ -d bin ] && rm -rf bin/*

.PHONY: build install test clean
