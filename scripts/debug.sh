#!/usr/bin/env bash
while true; do
    echo BUILDING...
#    go build -gcflags "all=-N -l" ./cmd/offensive/main.go
#    dlv exec ./main --headless --listen=:2345 --api-version=2 --log &
    dlv debug ./cmd/offensive/main.go --headless --accept-multiclient --listen=:2345 --api-version=2 --log
    echo Waiting for changes
    inotifywait -r -e modify -e move -e create -e delete ./internal ./cmd ./pkg
    pkill -P $$
done