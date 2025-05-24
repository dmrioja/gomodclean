package analyzer

import "fmt"

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
//   - ONLY_DIRECT: if the block contains only direct dependencies.
//   - ONLY_INDIRECT: if the block contains only indirect dependencies.
//   - MIXED: if the block contains both direct and indirect dependencies.
type consistency string

const (
	ONLY_DIRECT   consistency = "ONLY_DIRECT"
	ONLY_INDIRECT consistency = "ONLY_INDIRECT"
	MIXED         consistency = "MIXED"
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
			rb.consistency = ONLY_INDIRECT
		} else {
			rb.consistency = ONLY_DIRECT
		}
	case ONLY_INDIRECT:
		if !indirect {
			rb.consistency = MIXED
		}
	case ONLY_DIRECT:
		if indirect {
			rb.consistency = MIXED
		}
	default:
		// already mixed, nothing to do here!
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
func (rs *reqStmts) checkRule1() (issues []string) {
	if len(rs.directLines) > 1 {
		issues = append(issues, fmt.Sprintf("direct require lines should be grouped into blocks but found %d isolated require directives.", len(rs.directLines)))
	}
	if len(rs.indirectLines) > 1 {
		issues = append(issues, fmt.Sprintf("indirect require lines should be grouped into blocks but found %d isolated require directives.", len(rs.indirectLines)))
	}
	return
}

// checkRule2 asserts go.mod file only contains 2 require blocks.
func (rs *reqStmts) checkRule2() (issues []string) {
	if len(rs.blocks) > 2 {
		issues = append(issues, fmt.Sprintf("there should be a maximum of 2 require blocks but found %d.", len(rs.blocks)))
	}

	if len(issues) > 0 {
		return
	}

	// at this point there should be a maximum of 2 require blocks and
	// a maximum of 1 isolated direct or indirect line.
	for _, block := range rs.blocks {
		if block.consistency == ONLY_DIRECT && len(rs.directLines) > 0 {
			issues = append(issues, fmt.Sprintf("require directive \"%s %s\" should be inside block.", rs.directLines[0].name, rs.directLines[0].version))
		} else if block.consistency == ONLY_INDIRECT && len(rs.indirectLines) > 0 {
			issues = append(issues, fmt.Sprintf("require directive \"%s %s\" should be inside block.", rs.indirectLines[0].name, rs.indirectLines[0].version))
		}
	}

	return
}

// checkRule3 asserts the first require block only contains direct dependencies
// while the second one only contains indirect ones.
func (rs *reqStmts) checkRule3() (issues []string) {
	// rule #3.1: check the first require block only contains direct dependencies.
	if len(rs.blocks) > 0 {
		if rs.blocks[0].consistency != ONLY_DIRECT {
			issues = append(issues, "first require block should only contain direct dependencies")
		}
	}

	// rule #3.2: check the second require block only contains indirect dependencies.
	if len(rs.blocks) > 1 {
		if rs.blocks[1].consistency != ONLY_INDIRECT {
			issues = append(issues, "second require block should only contain indirect dependencies")
		}
	}

	return
}
