package main

import (
	"flag"
	"os"

	"github.com/magicsong/s2irun/pkg/run"
)

func main() {
	flag.Parse()
	os.Exit(run.App())
}
