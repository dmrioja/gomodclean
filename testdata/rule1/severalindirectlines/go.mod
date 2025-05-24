module github.com/dmrioja/gomodclean/testdata/rule1/severalindirectlines

go 1.24.2

require github.com/foo/foo v1.2.3

require github.com/bar/bar/v2 v2.0.0 // indirect
require github.com/cosa/cosita/v5 v5.3.3 // indirect