#!/usr/bin/env bash
while true; do
    echo BUILDING...
    go build ./cmd/offensive/main.go
    ./main &
    echo Waiting for changes
    inotifywait -r -e modify -e move -e create -e delete ./internal ./cmd
    pkill -P $$
done