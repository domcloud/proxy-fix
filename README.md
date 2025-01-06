# proxy-fix

This simple app simply spawns another HTTP server with `PORT` env injected with another number, then proxying it with [a clean header](https://github.com/domcloud/proxy-fix/blob/cb40ba29e5ac592438848b3071f2137ea9b3e0b6/main.go#L132-L140) request.

Built primarily for fixing https://github.com/phusion/passenger/issues/2521 temporarily, then added more stuff for DOM Cloud purposes. HTTP and Websocket is supported. 

## Install

Download from releases or build it and place it to `~/.local/bin/proxfix`.

```bash
PROXYFIX=proxy-fix-linux-$( [ "$(uname -m)" = "aarch64" ] && echo "arm64" || echo "amd64" )
wget https://github.com/domcloud/proxy-fix/releases/download/v0.2.5/$PROXYFIX.tar.gz
tar -xf $PROXYFIX.tar.gz && mv -f $PROXYFIX /usr/local/bin/proxfix && rm -rf $PROXYFIX*
```

## Usage

Use `Makefile` to build and run the app. Requires `make`, `go` and `bun` already installed.

```sh
make build
make run
```

## Envars

|Envar Example|Means|
|:--|:--|
|`NOHUP=1`|Run the app detached from this process tree, so when NGINX reloads it will not kill the process. Only works if used for only one instance per linux user.|
|`TARGET=127.0.0.1:1234`|Run the app with forwarding HTTP traffic to given IP and port instead of randomly. Useful for apps that can't customize the PORT via envar like Docker|

## Testing

Use `curl` and `wscat` to test with [test/app.ts](./test/app.ts).

```
curl -H '!~bad-headerz: x' -vvv localhost:8080
wscat -H '!~bad-headerz: x' -c "ws://localhost:8080/ws"
```
