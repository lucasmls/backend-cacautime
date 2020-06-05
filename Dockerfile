FROM golang:1.14.4-alpine3.11

RUN mkdir /app
WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY init.sh .

RUN apk add git
RUN apk add postgresql-client

RUN go mod download
COPY . .