FROM golang:1.19-alpine as builder

WORKDIR /go/src/github.com/kubesphere/s2irun
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY vendor/ vendor/
COPY go.mod go.sum ./

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o builder ./cmd/main.go

FROM alpine:3.11

WORKDIR /root/

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh
    
ENV S2I_CONFIG_PATH=/root/data/config.json
COPY --from=builder /go/src/github.com/kubesphere/s2irun/builder .
CMD ["./builder", "-v=4", "-logtostderr=true"]





