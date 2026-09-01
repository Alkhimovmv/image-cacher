.PHONY: build run test

build:
	mkdir -p app-data/images app-data/cache build
	go build -o build/image_cache ./cmd

run: build
	./build/image_cache

test:
	go test ./...
