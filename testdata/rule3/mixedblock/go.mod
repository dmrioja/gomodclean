module github.com/dmrioja/gomodclean/testdata/rule3/mixedblock

go 1.24.2

require (
    github.com/bar/bar/v2 v2.0.0
    github.com/dmrioja/shodo v1.0.0 // indirect
)

require (
    github.com/foo/foo v1.2.3
    github.com/cosa/cosita/v5 v5.3.3 // indirect
)
