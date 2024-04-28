FROM golang:1.22.2

WORKDIR /go/src/app

RUN go install github.com/cosmtrek/air@latest

COPY . .
RUN go mod tidy