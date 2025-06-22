package analyzer

import "fmt"

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

// reqLine represents a single require directive.
type reqLine struct {
	name     string
	version  string
	indirect bool
}

// reqBlock represents a single require block.
type reqBlock struct {
	lines       []*reqLine
	consistency consistency
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
func (rs *reqStmts) analyze() []string {
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

// checkRule1 asserts require lines and grouped into blocks.
//
// Note: if there is just one direct or indirect require directive
// there is no need to encapsulate it into a require block.
func (rs *reqStmts) checkRule1() []string {
	var issues []string

	if len(rs.directLines) > 1 {
		issues = append(issues, fmt.Sprintf("direct require lines should be grouped into blocks but found %d isolated require directives.", len(rs.directLines)))
	}

	if len(rs.indirectLines) > 1 {
		issues = append(issues, fmt.Sprintf("indirect require lines should be grouped into blocks but found %d isolated require directives.", len(rs.indirectLines)))
	}

	return issues
}

// checkRule2 asserts go.mod file only contains 2 require blocks.
func (rs *reqStmts) checkRule2() []string {
	var issues []string

	if len(rs.blocks) > maxRequireBlocks {
		issues = append(issues, fmt.Sprintf("there should be a maximum of 2 require blocks but found %d.", len(rs.blocks)))
	}

	if len(issues) > 0 {
		return issues
	}

	// at this point there should be a maximum of 2 require blocks and
	// a maximum of 1 isolated direct or indirect line.
	for _, block := range rs.blocks {
		if block.consistency == onlyDirect && len(rs.directLines) > 0 {
			issues = append(issues, fmt.Sprintf("require directive \"%s %s\" should be inside block.", rs.directLines[0].name, rs.directLines[0].version))
		} else if block.consistency == onlyIndirect && len(rs.indirectLines) > 0 {
			issues = append(issues, fmt.Sprintf("require directive \"%s %s\" should be inside block.", rs.indirectLines[0].name, rs.indirectLines[0].version))
		}
	}

	return issues
}

// checkRule3 asserts the first require block only contains direct dependencies
// while the second one only contains indirect ones.
func (rs *reqStmts) checkRule3() []string {
	var issues []string

	// case we have 2 require blocks
	// rule #3.1: check the first require block only contains direct dependencies.
	// rule #3.2: check the second require block only contains indirect dependencies.
	if len(rs.blocks) > 1 {
		if rs.blocks[0].consistency != onlyDirect {
			issues = append(issues, "first require block should only contain direct dependencies.")
		}

		if rs.blocks[1].consistency != onlyIndirect {
			issues = append(issues, "second require block should only contain indirect dependencies.")
		}

		if len(issues) > 0 {
			return issues
		}
	}

	// case we only have 1 require block
	if len(rs.blocks) > 0 {
		if rs.blocks[0].consistency == mixed {
			issues = append(issues, "first require block should not contain mixed dependencies.")
		}
	}

	return issues
}
