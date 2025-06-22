package main

import (
	"log"
	"os"

	"github.com/dmrioja/gomodclean/internal/analyzer"
	"golang.org/x/tools/go/analysis"
)

func main() {
	var pass *analysis.Pass

	issues, err := analyzer.AnalyzePass(pass)

	if err != nil {
		log.Fatal(err)
	}

	for _, issue := range issues {
		log.Println(issue)
	}

	os.Exit(0)
}
