package main

import (
	"os"

	"gomodclean/internal/analyzer"
)

func main() {
	os.Exit(analyzer.Run())
}
