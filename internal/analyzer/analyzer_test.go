package analyzer

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/mod/modfile"
)

func TestAnalyzeOneDirective(t *testing.T) {
	t.Parallel()

	file := retrieveGoModFile("rule1", "onedirective")

	issues := processFile(file).analyze()

	assert.Empty(t, issues)
}

func TestAnalyzeSeveralDirectLines(t *testing.T) {
	t.Parallel()

	file := retrieveGoModFile("rule1", "severaldirectlines")

	issues := processFile(file).analyze()

	assert.Len(t, issues, 1)
	assert.Equal(t, "direct require lines should be grouped into blocks but found 2 isolated require directives.", issues[0])
}

func TestAnalyzeSeveralIndirectLines(t *testing.T) {
	t.Parallel()

	file := retrieveGoModFile("rule1", "severalindirectlines")

	issues := processFile(file).analyze()

	assert.Len(t, issues, 1)
	assert.Equal(t, "indirect require lines should be grouped into blocks but found 2 isolated require directives.", issues[0])
}

func TestAnalyzeBothDirectAndIndirectLines(t *testing.T) {
	t.Parallel()

	file := retrieveGoModFile("rule1", "bothdirectandindirectlines")

	issues := processFile(file).analyze()

	assert.Len(t, issues, 2)
	assert.Equal(t, "direct require lines should be grouped into blocks but found 2 isolated require directives.", issues[0])
	assert.Equal(t, "indirect require lines should be grouped into blocks but found 2 isolated require directives.", issues[1])
}

func TestAnalyzeOnlyTwoRequireBlocks(t *testing.T) {
	t.Parallel()

	file := retrieveGoModFile("rule2", "onlytworequireblocks")

	issues := processFile(file).analyze()

	assert.Empty(t, issues)
}

func TestAnalyzeMoreThanTwoRequireBlocks(t *testing.T) {
	t.Parallel()

	file := retrieveGoModFile("rule2", "morethantworequireblocks")

	issues := processFile(file).analyze()

	assert.Len(t, issues, 1)
	assert.Equal(t, "there should be a maximum of 2 require blocks but found 4.", issues[0])
}

func TestAnalyzeIsolatedDirectLineShouldBeInsideBlock(t *testing.T) {
	t.Parallel()

	file := retrieveGoModFile("rule2", "isolateddirectlineshouldbeinsideblock")

	issues := processFile(file).analyze()

	assert.Len(t, issues, 1)
	assert.Equal(t, "require directive \"github.com/bar/bar/v2 v2.0.0\" should be inside block.", issues[0])
}

func TestAnalyzeIsolatedIndirectLineShouldBeInsideBlock(t *testing.T) {
	t.Parallel()

	file := retrieveGoModFile("rule2", "isolatedindirectlineshouldbeinsideblock")

	issues := processFile(file).analyze()

	assert.Len(t, issues, 1)
	assert.Equal(t, "require directive \"github.com/dmrioja/shodo v1.0.0\" should be inside block.", issues[0])
}

func TestAnalyzeIsolatedLinesShouldBeInsideBlock(t *testing.T) {
	t.Parallel()

	file := retrieveGoModFile("rule2", "isolatedlinesshouldbeinsideblock")

	issues := processFile(file).analyze()

	assert.Len(t, issues, 2)
	assert.Equal(t, "require directive \"github.com/bar/bar/v2 v2.0.0\" should be inside block.", issues[0])
	assert.Equal(t, "require directive \"github.com/cosa/cosita/v5 v5.3.3\" should be inside block.", issues[1])
}

func TestAnalyzeTwoCorrectBlocks(t *testing.T) {
	t.Parallel()

	file := retrieveGoModFile("rule3", "twocorrectblocks")

	issues := processFile(file).analyze()

	assert.Empty(t, issues)
}

func TestAnalyzeUnorderedBlocks(t *testing.T) {
	t.Parallel()

	file := retrieveGoModFile("rule3", "unorderedblocks")

	issues := processFile(file).analyze()

	assert.Len(t, issues, 2)
	assert.Equal(t, "first require block should only contain direct dependencies.", issues[0])
	assert.Equal(t, "second require block should only contain indirect dependencies.", issues[1])
}

func TestAnalyzeMixedBlock(t *testing.T) {
	t.Parallel()

	file := retrieveGoModFile("rule3", "mixedblock")

	issues := processFile(file).analyze()

	assert.Len(t, issues, 2)
	assert.Equal(t, "first require block should only contain direct dependencies.", issues[0])
	assert.Equal(t, "second require block should only contain indirect dependencies.", issues[1])
}

func TestAnalyzeOnlyOneIndirectBlock(t *testing.T) {
	t.Parallel()

	file := retrieveGoModFile("rule3", "onlyoneindirectblock")

	issues := processFile(file).analyze()

	assert.Empty(t, issues)
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
