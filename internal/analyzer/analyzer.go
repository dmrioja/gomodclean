package analyzer

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/mod/modfile"
)

func Run() int {

	// retrieve go.mod file
	file, err := getGoModFile()
	if err != nil {
		log.Fatal(err)
	}

	// process file (to extract the require statements)
	reqStmts := processFile(file)

	// analyze require staments
	issues := reqStmts.analyze()

	// analyze issues
	if len(issues) != 0 {
		for _, issue := range issues {
			fmt.Println(issue)
		}
		return 1
	}

	return 0
}

// processFile parses the go.mod file into a reqStmts struct.
func processFile(file *modfile.File) *reqStmts {
	reqStmts := &reqStmts{}

	for _, stmt := range file.Syntax.Stmt {
		switch _type := stmt.(type) {
		case *modfile.Line:
			if isRequire(_type.Token) {
				reqStmts.addLine(&reqLine{
					name:     _type.Token[1],
					version:  _type.Token[2],
					indirect: isIndirect(_type.Comment()),
				})
			}
		case *modfile.LineBlock:
			if isRequire(_type.Token) {
				block := &reqBlock{}
				for _, line := range _type.Line {
					// TODO: should we allow empty lines ??
					if len(line.Token) > 1 {
						block.addLine(&reqLine{
							name:     line.Token[0],
							version:  line.Token[1],
							indirect: isIndirect(line.Comment()),
						})
					}
				}
				reqStmts.addBlock(block)
			}
		default:
			// just do nothing
		}
	}

	return reqStmts
}

// isRequire returns true if the directive's first token is "require".
func isRequire(tokens []string) bool {
	return tokens[0] == "require"
}

// isIndirect returns true if the directive contains an "indirect" suffix.
func isIndirect(comment *modfile.Comments) bool {
	if comment != nil {
		for _, suffix := range comment.Suffix {
			if strings.Contains(suffix.Token, "indirect") {
				return true
			}
		}
	}
	return false
}
