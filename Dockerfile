FROM golang:1.11-alpine as builder

WORKDIR /go/src/github.com/magicsong/s2irun
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY vendor/ vendor/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o builder github.com/magicsong/s2irun/cmd

FROM ubuntu:latest

WORKDIR /root/
ENV S2I_CONFIG_PATH=/root/config.json
COPY --from=builder /go/src/github.com/magicsong/s2irun/builder .
ENTRYPOINT ["./builder"]





