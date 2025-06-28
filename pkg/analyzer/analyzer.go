package analyzer

import (
	"fmt"
	"go/token"

	"github.com/dmrioja/gomodclean/internal/io"
	"github.com/dmrioja/gomodclean/internal/processor"
)

// Issue represent an issue or problem detected during an analysis.
type Issue struct {
	Position token.Position
	Text     string
}

// Analyze run the analysis of gomodclean linter.
func Analyze() ([]Issue, error) {
	// retrieve go.mod file
	file, err := io.GetGoModFile()
	if err != nil {
		//nolint:err113
		return nil, fmt.Errorf("could not retrieve go.mod file: %s", err.Error())
	}

	issues := processor.ProcessFile(file)

	//nolint:prealloc
	var result []Issue
	for _, issue := range issues {
		result = append(result, Issue{
			Position: issue.Position,
			Text:     issue.Text,
		})
	}

	return result, nil
}
