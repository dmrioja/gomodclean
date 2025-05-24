# Rule 3

First require block must only contain direct dependencies while the second one must only contain indirect ones.

#### Bad:
```go
require (
    github.com/dmrioja/shodo v1.0.0 // indirect
    github.com/foo/foo v1.2.3
)

require (
    github.com/bar/bar/v2 v2.0.0
    github.com/cosa/cosita/v5 v5.3.3 // indirect
)
```

#### Good:
```go
require (
    github.com/bar/bar/v2 v2.0.0
    github.com/foo/foo v1.2.3
)

require (
    github.com/cosa/cosita/v5 v5.3.3 // indirect
    github.com/dmrioja/shodo v1.0.0 // indirect
)
```