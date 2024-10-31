# Build layer
FROM golang:latest AS build

WORKDIR /app

COPY . .

RUN go mod download

RUN GOOS=linux go build -o reaper

# Run layer
FROM ubuntu:latest

ENV DEBIAN_FRONTEND=noninteractive

RUN useradd -m -d /app -s /bin/bash app

WORKDIR /app
COPY . .
COPY --from=build /app/reaper .
RUN chown -R app /app
USER app

CMD ["./reaper"]
