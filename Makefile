.PHONY: build run

build:
	env go build -o ./build/proxy-fix ./...
	chmod +x ./build/proxy-fix

build-ci:
	env GOOS=linux GOARCH=amd64 go build -o ./build/proxy-fix-linux-amd64 -ldflags="-w -s" ./...
	env GOOS=linux GOARCH=arm64 go build -o ./build/proxy-fix-linux-arm64 -ldflags="-w -s" ./...
	cd ./build && tar -zcvf ./proxy-fix-linux-amd64.tar.gz ./proxy-fix-linux-amd64
	cd ./build && tar -zcvf ./proxy-fix-linux-arm64.tar.gz ./proxy-fix-linux-arm64

run:
	env PORT=8080 NOHUP=1 go run . bun ./test/app.ts
