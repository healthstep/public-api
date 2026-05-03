# syntax=docker/dockerfile:1
FROM golang:1.25-alpine AS build

RUN apk update && apk add --no-cache \
    git \
    gcc \
    musl-dev

ARG GITHUB_TOKEN
RUN echo "machine github.com login porebric password ${GITHUB_TOKEN}" > /root/.netrc && chmod 600 /root/.netrc

ENV GOPRIVATE=github.com

WORKDIR /app

COPY core-users/go.mod core-users/go.sum ./core-users/
COPY core-health/go.mod core-health/go.sum ./core-health/
COPY public-api/go.mod public-api/go.sum ./public-api/

WORKDIR /app/public-api
RUN go mod download

WORKDIR /app
COPY core-users/ ./core-users/
COPY core-health/ ./core-health/
COPY public-api/ ./public-api/

WORKDIR /app/public-api
RUN go mod tidy
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /out/public-api ./cmd/public-api

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=build /out/public-api .
COPY --from=build /app/public-api/config/configs_keys.yml ./config/configs_keys.yml
EXPOSE 8080
ENV APP_ENV=production
CMD ["./public-api"]
