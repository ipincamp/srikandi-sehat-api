FROM golang:1.24.6-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/migrate ./cmd/migrate/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/seeder ./cmd/seed/main.go

FROM alpine:latest

RUN apk add --no-cache tzdata bash

WORKDIR /app

COPY --from=build /app/server .
COPY --from=build /app/migrate .
COPY --from=build /app/seeder .

COPY serviceAccountKey.json .
COPY .docker/entrypoint.sh .

RUN chmod +x entrypoint.sh

EXPOSE 10000

ENTRYPOINT ["/app/entrypoint.sh"]
