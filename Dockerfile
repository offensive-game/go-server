FROM golang:1.11-alpine
LABEL maintainer="Djordje Vukovic"

RUN apk add --no-cache ca-certificates \
        dpkg \
        gcc \
        git \
        musl-dev \
        bash \
        inotify-tools

RUN go get github.com/derekparker/delve/cmd/dlv

WORKDIR /usr/src/app
COPY . .

RUN go mod tidy

RUN ["chmod", "+x", "./scripts/debug.sh"]
RUN ["chmod", "+x", "./scripts/run.sh"]

EXPOSE 7070
EXPOSE 2345