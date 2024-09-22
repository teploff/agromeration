.PHONY: build-local up-local

build-local:
	go build cmd/main.go

up-local:
	./main