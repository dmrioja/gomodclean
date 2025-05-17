package analyzer

import "fmt"

// reqStmts stands for "require statements". It holds the full
// require directives encapsulated into lines and blocks.
type reqStmts struct {

	// lines contains all the isolated "require" directive lines.
	lines []*reqLine

	// blocks contains all the "require" directive blocks.
	blocks []*reqBlock
}

// reqLine represents a single require directive.
type reqLine struct {

	// TODO: name and version could be useless. I used them for the first iteration to print
	// the dependency but I don't think a need them any more.
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
	rs.lines = append(rs.lines, line)
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
func (rs *reqStmts) analyze() (issues []string) {

	// rule #1: check require lines are grouped into blocks.
	if len(rs.lines) > 1 {
		// TODO: there could be 1 direct and 1 indirect line
		issues = append(issues, fmt.Sprintf("require lines should be grouped into blocks but found %d isolated require directives.", len(rs.lines)))
	} else if len(rs.lines) == 1 {
		// TODO: there could be 1 direct line and 1 indirect block or...
		// there could be 1 indirect line and 1 direct block
	}

	// rule #2: check go.mod file only contains 2 require blocks.
	if len(rs.blocks) > 2 {
		issues = append(issues, fmt.Sprintf("there should be a maximum of 2 require blocks but found %d", len(rs.blocks)))
	}

	// rule #3.1: check the first require block only contains direct dependencies.
	if len(rs.blocks) > 0 {
		if rs.blocks[0].consistency != ONLY_DIRECT {
			issues = append(issues, "first require block should only contain direct dependencies")
		}
	}

	// rule #3.2: check the second require block only contains indirect dependencies.
	if len(rs.blocks) > 1 {
		if rs.blocks[0].consistency != ONLY_INDIRECT {
			issues = append(issues, "second require block should only contain indirect dependencies")
		}
	}

	return
}
