  
# STAGE 0: Contruct build base
FROM golang:1.14-stretch as builder_base

WORKDIR /cacautime

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

COPY go.mod .
COPY go.sum .

RUN go mod download

# STAGE 1: Build binaries
FROM builder_base as builder
WORKDIR /cacautime
COPY . /cacautime

RUN go build -a -installsuffix cgo -o server github.com/lucasmls/backend-cacautime/cmd/server

# STAGE 2: Build server
FROM alpine as server
ENV PORT=3000
COPY --from=builder /cacautime/server /go/bin/server
RUN apk add -U --no-cache ca-certificates
EXPOSE 3000
ENTRYPOINT /go/bin/server