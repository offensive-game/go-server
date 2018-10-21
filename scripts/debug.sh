#!/usr/bin/env bash
#go build -gcflags "all=-N -l" ./cmd/offensive/main.go
#dlv exec ./main --headless --listen=:2345 --api-version=2

dlv debug ./cmd/offensive/main.go -l 0.0.0.0:2345 --headless=true --log=true