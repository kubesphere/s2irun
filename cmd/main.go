package main

import (
	"flag"
	"os"

	"github.com/MagicSong/s2irun/pkg/run"
)

func main() {
	flag.Parse()
	os.Exit(run.App())
}
