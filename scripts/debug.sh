#!/usr/bin/env sh
dlv debug ./simple-http.go --headless --accept-multiclient --listen=:2345 --api-version=2 --log