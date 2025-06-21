package analyzer

import (
	"github.com/dmrioja/gomodclean/internal/analyzer"

	"golang.org/x/tools/go/analysis"
)

func Run(pass *analysis.Pass) ([]string, error) {

	// retrieve go.mod file
	file, err := analyzer.GetGoModFile()
	if err != nil {
		return []string{}, err
	}

	// process file (to extract the require statements)
	reqStmts := analyzer.ProcessFile(file)

	// analyze require staments
	return reqStmts.Analyze(), nil
}
