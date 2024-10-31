# Build layer
FROM golang:latest AS build

WORKDIR /app

COPY . .

RUN go mod download

ENV GOOS=linux

RUN go build -ldflags="-s -w" -o reaper ./cmd/reaper

# Run layer
FROM ubuntu:latest

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates

RUN useradd -m -d /app -s /bin/bash app

WORKDIR /app
COPY . .
COPY --from=build /app/reaper .
RUN chown -R app /app
USER app

CMD ["./reaper"]
