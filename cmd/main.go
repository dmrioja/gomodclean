package main

import (
	"log"
	"os"

	"github.com/dmrioja/gomodclean/pkg/analyzer"
)

func main() {
	issues, err := analyzer.Analyze()

	if err != nil {
		log.Fatal(err)
	}

	for _, issue := range issues {
		log.Println(issue)
	}

	os.Exit(0)
}
