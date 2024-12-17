build:
	[ -d bin ] || mkdir bin
	go build -o bin .

.PHONY: build
