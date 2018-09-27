#!/usr/bin/env bash
while true; do
    echo BUILDING...
    ps -A
    go build ./cmd/offensive/main.go
    ./main &
    echo Waiting for changes
    inotifywait -r -e modify -e move -e create -e delete ./internal ./cmd ./pkg
    pkill -P $$
done