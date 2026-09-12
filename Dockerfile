# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/medical-visit-api ./cmd/api

FROM alpine:3.21

RUN apk add --no-cache ca-certificates && addgroup -S app && adduser -S -G app app

WORKDIR /app
COPY --from=builder /out/medical-visit-api ./medical-visit-api
COPY migrations ./migrations

USER app
EXPOSE 8080

ENTRYPOINT ["./medical-visit-api"]
