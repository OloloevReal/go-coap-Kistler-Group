FROM golang:1.13.1-alpine3.10 AS build
RUN apk add --no-cache curl git build-base
WORKDIR $GOPATH/src/github.com/OloloevReal/go-coap-Kistler-Group
COPY go.mod go.sum ./
RUN go mod download
COPY . .

