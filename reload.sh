#!/usr/bin/env bash

set -e

echo "replace api and websocket ip address..."
addr_http=${HOST_IP}:7070
addr_ws=${HOST_IP}:7000

echo "build gochat.bin ..."
# CGO_CFLAGS="-g -O2 -Wno-return-local-addr" fix compile sqlite3 warning
CGO_CFLAGS="-g -O2 -Wno-return-local-addr" go build -o /go/src/gochat/bin/gochat.bin /go/src/gochat/main.go
echo "restart all ..."
supervisorctl restart all
echo "all Done."
echo "Beautiful ! start the world ! "
