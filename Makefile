build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o _output/cmd/builder github.com/MagicSong/s2irun/cmd
run: 
	go run ./cmd/main.go