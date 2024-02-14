.PHONY: build run

build:
	env go build -o ./build/proxy-fix ./...
	chmod +x ./build/proxy-fix

build-ci:
	env GOOS=linux GOARCH=amd64 go build -o ./build/proxy-fix-linux-amd64 ./...
	env GOOS=linux GOARCH=arm64 go build -o ./build/proxy-fix-linux-arm64 ./...

run:
	env PORT=8080 go run . bun ./test/app.ts
