# gomodclean

Linter to check dependencies are well structured inside your go.mod file.

<div>
    <img src="docs/gopher.png" alt="gomodclean gopher logo" width="400"/>
</div>

## Details

### #1: Check require lines are grouped into blocks

<table style="witdh:100%">
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody><tr><td>

```go
require github.com/bar/bar v2.0.0
require github.com/foo/foo v1.2.3
```

</td><td>

```go
require (
    github.com/bar/bar v2.0.0
    github.com/foo/foo v1.2.3
)
```

</td></tr></tbody>
</table>

#### Note:
If there is just one direct or indirect require directive, there is no need to encapsulate it into a require block, so the following example is valid:

```go
require github.com/foo/foo v1.2.3

require (
    github.com/bar/bar v2.0.0 // indirect
    github.com/cosa/cosita v5.3.3 // indirect
)
```

### #2: Check go.mod file only contains 2 require blocks

<table style="witdh:100%">
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody><tr><td>

```go
require (
    github.com/foo/foo v1.2.3
)

require (
    github.com/bar/bar v2.0.0
)

require (
    github.com/cosa/cosita v5.3.3 // indirect
)
```

</td><td>

```go
require (
    github.com/bar/bar v2.0.0
    github.com/foo/foo v1.2.3
)

require (
    github.com/cosa/cosita v5.3.3 // indirect
)
```

</td></tr></tbody>
</table>

### #3: Check the first require block only contains direct dependencies while the second one only contains indirect ones

<table style="witdh:100%">
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody><tr><td>

```go
require (
    github.com/dmrioja/shodo v1.0.0 // indirect
    github.com/foo/foo v1.2.3
)

require (
    github.com/bar/bar v2.0.0
    github.com/cosa/cosita v5.3.3 // indirect
)
```

</td><td>

```go
require (
    github.com/bar/bar v2.0.0
    github.com/foo/foo v1.2.3
)

require (
    github.com/cosa/cosita v5.3.3 // indirect
    github.com/dmrioja/shodo v1.0.0 // indirect
)
```

</td></tr></tbody>
</table>

## Usage

🚧 Work in Progress...