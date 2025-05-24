# Rule 2

go.mod file can only contain a maximum of 2 require blocks.

#### Bad:
```go
require (
    github.com/foo/foo v1.2.3
)

require (
    github.com/bar/bar/v2 v2.0.0
)

require (
    github.com/cosa/cosita/v5 v5.3.3 // indirect
)

require (
    github.com/dmrioja/shodo v1.0.0 // indirect
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