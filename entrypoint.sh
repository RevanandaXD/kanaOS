#!/bin/sh
set -eu

: "${PORT:=8080}"
export PORT

PORT=3000 HOST=127.0.0.1 node /app/web/build/index.js &
web_pid=$!

kanaos &
api_pid=$!

trap 'kill -TERM "$web_pid" "$api_pid" 2>/dev/null || true' TERM INT

while kill -0 "$web_pid" 2>/dev/null && kill -0 "$api_pid" 2>/dev/null; do
  sleep 1
done

kill -TERM "$web_pid" "$api_pid" 2>/dev/null || true
wait "$web_pid" 2>/dev/null || true
wait "$api_pid" 2>/dev/null || true
