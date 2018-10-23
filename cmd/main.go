package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/MagicSong/s2irun/pkg/api"
	"github.com/MagicSong/s2irun/pkg/run"
)

func main() {
	flag.Parse()
	apiConfig := new(api.Config)
	err := run.S2I(apiConfig)
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
	fmt.Println("Hello World")
}
