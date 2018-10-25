FROM golang:1.11-alpine as builder

WORKDIR /go/src/github.com/MagicSong/s2irun
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY vendor/ vendor/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o builder github.com/MagicSong/s2irun/cmd

FROM alpine:latest

WORKDIR /root/

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh
    
ENV S2I_CONFIG_PATH=/root/data/config.json
COPY --from=builder /go/src/github.com/MagicSong/s2irun/builder .
CMD ["./builder"]





