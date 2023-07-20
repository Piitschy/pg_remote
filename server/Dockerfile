# syntax=docker/dockerfile:1

FROM golang:1.20-alpine

RUN apk add postgresql

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY ./static ./static

RUN go build -o /bin

EXPOSE 3000

CMD [ "/bin" ]