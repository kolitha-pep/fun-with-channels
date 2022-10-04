FROM golang:1.18-bullseye

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .