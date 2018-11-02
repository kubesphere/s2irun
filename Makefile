IMG ?= magicsong/s2irunner:latest
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o _output/cmd/builder github.com/magicsong/s2irun/cmd
run: 
	S2I_CONFIG_PATH=test/config.json go run ./cmd/main.go
image:
	docker build . -t ${IMG}