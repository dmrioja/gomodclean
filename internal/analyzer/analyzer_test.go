package analyzer

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/mod/modfile"
)

// testCase represent a case inside a rule test.
type testCase struct {
	name           string
	expectedIssues []string
}

//nolint:gochecknoglobals
var scenarios = []struct {
	rule      string
	testCases []testCase
}{
	{
		rule: "rule1",
		testCases: []testCase{
			{name: "onedirective", expectedIssues: nil},
			{name: "severaldirectlines", expectedIssues: []string{
				"direct require lines should be grouped into blocks but found 2 isolated require directives.",
			}},
			{name: "severalindirectlines", expectedIssues: []string{
				"indirect require lines should be grouped into blocks but found 2 isolated require directives.",
			}},
			{name: "bothdirectandindirectlines", expectedIssues: []string{
				"direct require lines should be grouped into blocks but found 2 isolated require directives.",
				"indirect require lines should be grouped into blocks but found 2 isolated require directives.",
			}},
		},
	},
	{
		rule: "rule2",
		testCases: []testCase{
			{name: "onlytworequireblocks", expectedIssues: nil},
			{name: "morethantworequireblocks", expectedIssues: []string{
				"there should be a maximum of 2 require blocks but found 4.",
			}},
			{name: "isolateddirectlineshouldbeinsideblock", expectedIssues: []string{
				"require directive \"github.com/bar/bar/v2 v2.0.0\" should be inside block.",
			}},
			{name: "isolatedindirectlineshouldbeinsideblock", expectedIssues: []string{
				"require directive \"github.com/dmrioja/shodo v1.0.0\" should be inside block.",
			}},
			{name: "isolatedlinesshouldbeinsideblock", expectedIssues: []string{
				"require directive \"github.com/bar/bar/v2 v2.0.0\" should be inside block.",
				"require directive \"github.com/cosa/cosita/v5 v5.3.3\" should be inside block.",
			}},
		},
	},
	{
		rule: "rule3",
		testCases: []testCase{
			{name: "twocorrectblocks", expectedIssues: nil},
			{name: "unorderedblocks", expectedIssues: []string{
				"first require block should only contain direct dependencies.",
				"second require block should only contain indirect dependencies.",
			}},
			{name: "mixedblock", expectedIssues: []string{
				"first require block should only contain direct dependencies.",
				"second require block should only contain indirect dependencies.",
			}},
			{name: "onlyoneindirectblock", expectedIssues: nil},
			{name: "indirectcomment", expectedIssues: nil},
		},
	},
}

func TestAnalyzeScenarios(t *testing.T) {
	t.Parallel()

	for _, scenario := range scenarios {
		t.Run(scenario.rule, func(t *testing.T) {
			t.Parallel()

			for _, testCase := range scenario.testCases {
				t.Run(testCase.name, func(t *testing.T) {
					t.Parallel()

					file := retrieveGoModFile(scenario.rule, testCase.name)
					issues := processFile(file).analyze()

					if testCase.expectedIssues == nil {
						assert.Empty(t, issues)
					} else {
						assert.Len(t, issues, len(testCase.expectedIssues))

						for i, expectedIssue := range testCase.expectedIssues {
							assert.Equal(t, expectedIssue, issues[i])
						}
					}
				})
			}
		})
	}
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
