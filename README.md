# proxy-fix

This simple app simply spawns another HTTP server with `PORT` env injected with another number, then proxying it with [a clean header](https://github.com/domcloud/proxy-fix/blob/cb40ba29e5ac592438848b3071f2137ea9b3e0b6/main.go#L132-L140) request.

Built primarily for fixing https://github.com/phusion/passenger/issues/2521 temporarily. HTTP and Websocket is supported.

## Install

Download from releases or build it and place it to `~/.local/bin/bunfix`

## Usage

Use `Makefile` to build and run the app. Requires `make`, `go` and `bun` already installed.

```sh
make build
make run
```

## Testing

Use `curl` and `wscat` to test with [test/app.ts](./test/app.ts).

```
curl -H '!~bad-headerz: x' -vvv localhost:8080
wscat -H '!~bad-headerz: x' -c "ws://localhost:8080/ws"
```
