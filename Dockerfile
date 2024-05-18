# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS build-stage
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY ./cmd/server/ ./cmd/server/
COPY ./internal/ ./internal/
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN cd /app/cmd/server/ && swag init
RUN cd /app && go build ./cmd/server

FROM build-stagen AS test-stage
RUN go test ./...

FROM alpine AS production-stage
RUN apk add postgresql
WORKDIR /
COPY --from=build-stage /app/server /server
RUN chmod +x /server

EXPOSE 3000

CMD [ "/server" ]
