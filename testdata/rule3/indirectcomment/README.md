# indirectcomment

The go.mod file is valid since:

- It only contains 2 required blocks.
- There are no isolated direct lines.
- The first block only contains direct lines.
- The second block only contains indirect lines.
- One of the direct dependencies contains a comment which contains the string "// indirect" (but it is not an indirect dependency so the analysis should pass ok).

