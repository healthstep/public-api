FROM golang:1.25-alpine AS builder
WORKDIR /build

RUN apk add --no-cache git

COPY core-users/go.mod core-users/go.sum ./core-users/
COPY core-health/go.mod core-health/go.sum ./core-health/
COPY public-api/go.mod public-api/go.sum ./public-api/

WORKDIR /build/public-api
RUN go mod download

WORKDIR /build
COPY core-users/ ./core-users/
COPY core-health/ ./core-health/
COPY public-api/ ./public-api/

WORKDIR /build/public-api
RUN CGO_ENABLED=0 go build -o /app/public-api ./cmd/public-api

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/public-api .
COPY --from=builder /build/public-api/config/configs_keys.yml ./config/configs_keys.yml
EXPOSE 8080
CMD ["./public-api"]
