package analyzer

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/mod/modfile"
)

func TestAnalyze(t *testing.T) {
	file, err := readGoModFile("../../testdata/rule1/onedirective/go.mod")
	if err != nil {
		t.Fatal(err)
	}

	issues := processFile(file).analyze()

	assert.Len(t, issues, 0)
}

func readGoModFile(filepath string) (*modfile.File, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("could not read go.mod file: %w", err)
	}

	file, err := modfile.Parse("go.mod", content, nil)
	if err != nil {
		return nil, fmt.Errorf("could not parse go.mod file: %w", err)
	}

	return file, nil
}
