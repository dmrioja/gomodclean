package processor

import (
	"fmt"
	"go/token"
	"strings"

	"golang.org/x/mod/modfile"
)

const (
	maxRequireBlocks = 2
)

// reqStmts stands for "require statements". It holds the full
// require directives encapsulated into lines and blocks.
type reqStmts struct {

	// directLines contains all the isolated direct "require" directive lines.
	directLines []*reqLine

	// indirectLines contains all the isolated indirect "require" directive lines.
	indirectLines []*reqLine

	// blocks contains all the "require" directive blocks.
	blocks []*reqBlock
}

type Issue struct {
	Position token.Position
	Text     string
}

// reqLine represents a single require directive.
type reqLine struct {
	name     string
	version  string
	indirect bool
	position token.Position
}

// reqBlock represents a single require block.
type reqBlock struct {
	lines       []*reqLine
	consistency consistency
	position    token.Position
}

// consistency is an enum representing all the possible cases for the
// collection of dependencies a block could contain. This cases are:
//   - onlyDirect: if the block contains only direct dependencies.
//   - onlyIndirect: if the block contains only indirect dependencies.
//   - mixed: if the block contains both direct and indirect dependencies.
type consistency string

const (
	onlyDirect   consistency = "ONLY_DIRECT"
	onlyIndirect consistency = "ONLY_INDIRECT"
	mixed        consistency = "MIXED"
)

// addLine adds a new line to the reqStmts isolated require directive lines.
func (rs *reqStmts) addLine(line *reqLine) {
	if line.indirect {
		rs.indirectLines = append(rs.indirectLines, line)
	} else {
		rs.directLines = append(rs.directLines, line)
	}
}

// addBlock adds a whole block of require directive lines to the reqStmts blocks.
func (rs *reqStmts) addBlock(block *reqBlock) {
	rs.blocks = append(rs.blocks, block)
}

// addLine adds a new require directive line to the block.
// It also updates the consistency of the block.
func (rb *reqBlock) addLine(line *reqLine) {
	rb.lines = append(rb.lines, line)
	rb.updateConsistency(line.indirect)
}

// updateConsistency updates the block's consistency according to
// the new added require directive line.
func (rb *reqBlock) updateConsistency(indirect bool) {
	switch rb.consistency {
	case "":
		if indirect {
			rb.consistency = onlyIndirect
		} else {
			rb.consistency = onlyDirect
		}
	case onlyIndirect:
		if !indirect {
			rb.consistency = mixed
		}
	case onlyDirect:
		if indirect {
			rb.consistency = mixed
		}
	default:
		// already mixed, nothing to do here!
		return
	}
}

// analyze checks go.mod file satisfies all the rules defined for gomodclean.
func (rs *reqStmts) analyze() []Issue {
	// rule #1: check require lines are grouped into blocks.
	if issues := rs.checkRule1(); len(issues) > 0 {
		return issues
	}

	// rule #2: check go.mod file only contains 2 require blocks.
	if issues := rs.checkRule2(); len(issues) > 0 {
		return issues
	}

	// rule #3: check the first require block only contains direct dependencies
	// while the second one only contains indirect ones.
	return rs.checkRule3()
}

// checkRule1 asserts require lines are grouped into blocks.
//
// Note: if there is just one direct or indirect require directive
// there is no need to encapsulate it into a require block.
func (rs *reqStmts) checkRule1() []Issue {
	var issues []Issue

	if len(rs.directLines) > 1 {
		issues = append(issues, Issue{
			Position: rs.directLines[0].position,
			Text:     fmt.Sprintf("direct require lines should be grouped into blocks but found %d isolated require directives.", len(rs.directLines)),
		})
	}

	if len(rs.indirectLines) > 1 {
		issues = append(issues, Issue{
			Position: rs.indirectLines[0].position,
			Text:     fmt.Sprintf("indirect require lines should be grouped into blocks but found %d isolated require directives.", len(rs.indirectLines)),
		})
	}

	return issues
}

// checkRule2 asserts go.mod file only contains 2 require blocks.
func (rs *reqStmts) checkRule2() []Issue {
	var issues []Issue

	if len(rs.blocks) > maxRequireBlocks {
		issues = append(issues, Issue{
			Position: rs.blocks[0].position,
			Text:     fmt.Sprintf("there should be a maximum of 2 require blocks but found %d.", len(rs.blocks)),
		})

		return issues
	}

	// at this point there should be a maximum of 2 require blocks and
	// a maximum of 1 isolated direct or indirect line.
	for _, block := range rs.blocks {
		if block.consistency == onlyDirect && len(rs.directLines) > 0 {
			issues = append(issues, Issue{
				Position: rs.directLines[0].position,
				Text:     fmt.Sprintf("require directive \"%s %s\" should be inside block.", rs.directLines[0].name, rs.directLines[0].version),
			})
		} else if block.consistency == onlyIndirect && len(rs.indirectLines) > 0 {
			issues = append(issues, Issue{
				Position: rs.indirectLines[0].position,
				Text:     fmt.Sprintf("require directive \"%s %s\" should be inside block.", rs.indirectLines[0].name, rs.indirectLines[0].version),
			})
		}
	}

	return issues
}

// checkRule3 asserts the first require block only contains direct dependencies
// while the second one only contains indirect ones.
func (rs *reqStmts) checkRule3() []Issue {
	var issues []Issue

	// case we have 2 require blocks
	// rule #3.1: check the first require block only contains direct dependencies.
	// rule #3.2: check the second require block only contains indirect dependencies.
	if len(rs.blocks) > 1 {
		if rs.blocks[0].consistency != onlyDirect {
			issues = append(issues, Issue{
				Position: rs.blocks[0].position,
				Text:     "first require block should only contain direct dependencies.",
			})
		}

		if rs.blocks[1].consistency != onlyIndirect {
			issues = append(issues, Issue{
				Position: rs.blocks[1].position,
				Text:     "second require block should only contain indirect dependencies.",
			})
		}

		if len(issues) > 0 {
			return issues
		}
	}

	// case we only have 1 require block
	if len(rs.blocks) > 0 {
		if rs.blocks[0].consistency == mixed {
			issues = append(issues, Issue{
				Position: rs.blocks[0].position,
				Text:     "require block should not contain mixed dependencies.",
			})
		}
	}

	return issues
}

// ProcessFile parses the go.mod file into a reqStmts struct.
func ProcessFile(file *modfile.File) []Issue {
	reqStmts := &reqStmts{}

	for _, stmt := range file.Syntax.Stmt {
		switch _type := stmt.(type) {
		case *modfile.Line:
			if isRequire(_type.Token) {
				reqStmts.addLine(&reqLine{
					name:     _type.Token[1],
					version:  _type.Token[2],
					indirect: isIndirect(_type.Comment()),
					position: token.Position{
						Filename: file.Syntax.Name,
						Line:     _type.Start.Line,
					},
				})
			}
		case *modfile.LineBlock:
			if isRequire(_type.Token) {
				block := &reqBlock{
					position: token.Position{
						Filename: file.Syntax.Name,
						Line:     _type.Start.Line,
					},
				}

				for _, line := range _type.Line {
					// TODO: should we allow empty lines ??
					if len(line.Token) > 1 {
						block.addLine(&reqLine{
							name:     line.Token[0],
							version:  line.Token[1],
							indirect: isIndirect(line.Comment()),
							position: token.Position{
								Line: line.Start.Line,
							},
						})
					}
				}

				reqStmts.addBlock(block)
			}
		default:
			// just do nothing
			continue
		}
	}

	return reqStmts.analyze()
}

// isRequire returns true if the directive's first token is "require".
func isRequire(tokens []string) bool {
	return tokens[0] == "require"
}

// isIndirect returns true if the directive contains an "indirect" suffix.
func isIndirect(comment *modfile.Comments) bool {
	if comment != nil {
		for _, suffix := range comment.Suffix {
			if strings.HasPrefix(suffix.Token, "// indirect") {
				return true
			}
		}
	}

	return false
}
