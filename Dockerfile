FROM golang:1.13.0
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o main ./cmd/server
CMD ["/app/cmd/server/main"]
