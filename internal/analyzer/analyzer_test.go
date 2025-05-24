package analyzer

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/mod/modfile"
)

func TestAnalyzeOneDirective(t *testing.T) {
	file := retrieveGoModFile("rule1", "onedirective")

	issues := processFile(file).analyze()

	assert.Len(t, issues, 0)
}

func TestAnalyzeSeveralDirectLines(t *testing.T) {
	file := retrieveGoModFile("rule1", "severaldirectlines")

	issues := processFile(file).analyze()

	assert.Len(t, issues, 1)
	assert.Equal(t, issues[0], "direct require lines should be grouped into blocks but found 2 isolated require directives.")
}

func TestAnalyzeSeveralIndirectLines(t *testing.T) {
	file := retrieveGoModFile("rule1", "severalindirectlines")

	issues := processFile(file).analyze()

	assert.Len(t, issues, 1)
	assert.Equal(t, issues[0], "indirect require lines should be grouped into blocks but found 2 isolated require directives.")
}

func TestAnalyzeBothDirectAndIndirectLines(t *testing.T) {
	file := retrieveGoModFile("rule1", "bothdirectandindirectlines")

	issues := processFile(file).analyze()

	assert.Len(t, issues, 2)
	assert.Equal(t, issues[0], "direct require lines should be grouped into blocks but found 2 isolated require directives.")
	assert.Equal(t, issues[1], "indirect require lines should be grouped into blocks but found 2 isolated require directives.")
}

func TestAnalyzeOnlyTwoRequireBlocks(t *testing.T) {
	file := retrieveGoModFile("rule2", "onlytworequireblocks")

	issues := processFile(file).analyze()

	assert.Len(t, issues, 0)
}

func TestAnalyzeMoreThanTwoRequireBlocks(t *testing.T) {
	file := retrieveGoModFile("rule2", "morethantworequireblocks")

	issues := processFile(file).analyze()

	assert.Len(t, issues, 1)
	assert.Equal(t, issues[0], "there should be a maximum of 2 require blocks but found 4.")
}

func TestAnalyzeIsolatedDirectLineShouldBeInsideBlock(t *testing.T) {
	file := retrieveGoModFile("rule2", "isolateddirectlineshouldbeinsideblock")

	issues := processFile(file).analyze()

	assert.Len(t, issues, 1)
	assert.Equal(t, issues[0], "require directive \"github.com/bar/bar/v2 v2.0.0\" should be inside block.")
}

func TestAnalyzeIsolatedIndirectLineShouldBeInsideBlock(t *testing.T) {
	file := retrieveGoModFile("rule2", "isolatedindirectlineshouldbeinsideblock")

	issues := processFile(file).analyze()

	assert.Len(t, issues, 1)
	assert.Equal(t, issues[0], "require directive \"github.com/dmrioja/shodo v1.0.0\" should be inside block.")
}

func retrieveGoModFile(rule, testCase string) *modfile.File {
	file, err := readGoModFile(fmt.Sprintf("../../testdata/%s/%s/go.mod", rule, testCase))
	if err != nil {
		panic(err)
	}
	return file
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
