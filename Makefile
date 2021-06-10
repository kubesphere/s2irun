IMG ?= kubesphere/s2irun:v0.0.1
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build cmd
run:
	S2I_CONFIG_PATH=test/config.json go run ./cmd/main.go
run-b2i:
	S2I_CONFIG_PATH=test/b2iconfig.json go run ./cmd/main.go
image:
	docker build . -t ${IMG}
push: image
	docker push ${IMG}
test: fmt vet
	go test -v ./pkg/...
fmt:
	go fmt ./pkg/... ./cmd/...
vet:
	go vet ./pkg/... ./cmd/...