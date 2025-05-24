module github.com/dmrioja/gomodclean/testdata/rule2/isolatedindirectlineshouldbeinsideblock

go 1.24.2

require github.com/bar/bar/v2 v2.0.0

require (
    github.com/foo/foo v1.2.3 // indirect
    github.com/cosa/cosita/v5 v5.3.3 // indirect
)

require github.com/dmrioja/shodo v1.0.0 // indirect
