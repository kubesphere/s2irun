FROM golang:1.23 as builder

WORKDIR /go/src/github.com/kubesphere/s2irun
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY go.mod go.mod
COPY go.sum go.sum

# Build
RUN CGO_ENABLED=0 GO111MODULE=on GOOS=linux go build -a -o builder github.com/kubesphere/s2irun/cmd

FROM alpine:3.21

WORKDIR /root/

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh
    
ENV S2I_CONFIG_PATH=/root/data/config.json
COPY --from=builder /go/src/github.com/kubesphere/s2irun/builder .
CMD ["./builder", "-v=4", "-logtostderr=true"]





