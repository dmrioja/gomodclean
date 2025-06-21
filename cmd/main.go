package main

import (
	"os"

	"github.com/dmrioja/gomodclean/internal/analyzer"
)

func main() {
	os.Exit(analyzer.Run())
}
