# syntax=docker/dockerfile:1

FROM golang:1.20-alpine

RUN apk add postgresql

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init
RUN go build -o ./bin
RUN chmod +x ./bin

EXPOSE 3000

CMD [ "./bin" ]