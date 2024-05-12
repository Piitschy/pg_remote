# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS build-stage

RUN apk add postgresql

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./cmd/server/*.go ./
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init
RUN go build -o /server

FROM build-stagen AS test-stage
RUN go test ./...

FROM alpine AS production-stage
WORKDIR /
COPY --from=build-stage /server /server
RUN chmod +x /server

EXPOSE 3000

CMD [ "/server" ]
