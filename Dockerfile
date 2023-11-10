FROM golang:1.21-alpine AS builder
RUN apk --update add ca-certificates make npm
WORKDIR /build

COPY . /build

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN make build
RUN go build -ldflags="-s -w" -o reaper ./cmd/reaper

FROM scratch
COPY --from=builder ["/etc/ssl/certs/ca-certificates.crt", "/etc/ssl/certs/"]
COPY --from=builder ["/build/reaper", "/"]

ENTRYPOINT ["/reaper"]