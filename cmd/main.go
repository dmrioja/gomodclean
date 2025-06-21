package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dmrioja/gomodclean/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

func main() {
	var pass *analysis.Pass
	issues, err := analyzer.Run(pass)
	if err != nil {
		log.Fatal(err)
	}

	for _, issue := range issues {
		fmt.Println(issue)
	}

	os.Exit(0)
}
