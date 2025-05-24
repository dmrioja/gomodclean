# Rule 1

Require lines must be grouped into blocks.

#### Bad:
```go
require github.com/bar/bar/v2 v2.0.0
require github.com/foo/foo v1.2.3
```

#### Good:
```go
require (
    github.com/bar/bar/v2 v2.0.0
    github.com/foo/foo v1.2.3
)
```

**Note**: \
If there is just one direct or one indirect require directive there is no need to encapsulate it into a require block, so the following example is valid:
```go
require github.com/foo/foo v1.2.3

require (
    github.com/bar/bar/v2 v2.0.0 // indirect
    github.com/cosa/cosita/v5 v5.3.3 // indirect
)
```