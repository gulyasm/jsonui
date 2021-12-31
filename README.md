# JSONUI
[![](https://travis-ci.org/gulyasm/jsonui.svg?branch=master)](https://travis-ci.org/gulyasm/jsonui) [![](https://goreportcard.com/badge/github.com/gulyasm/jsonui)](https://goreportcard.com/report/github.com/gulyasm/jsonui)

`jsonui` is an interactive JSON explorer in your command line. You can pipe any JSON into `jsonui` and explore it, copy the path for each element.

![](img/jsonui.gif)

## Install
`go install github.com/gulyasm/jsonui@latest`
Binary will be installed into `${GOPATH}/bin`, where `GOPATH` defaults to `~/go`. Make sure the bin directory is in your `$PATH`.

## Binary Releases
[Binary releases are availabe](https://github.com/gulyasm/jsonui/releases)

## Usage
Just use the standard output:
```
cat test_big.json | jsonui
```

### Keys

#### `j`, `DownArrow`
Move down a line

#### `k`, `DownUp`
Move up a line

#### `J/PageDown`
Move down 15 lines

#### `K/PageUp`
Move up 15 lines

#### `h/?`
Toggle Help view

#### `e`
Toggle node (expend or collapse)

#### `E`
Expand all nodes

#### `C`
Collapse all nodes

#### `q/Ctrl+C`
Quit jsonui


## Acknowledgments
Special thanks for [asciimoo](https://github.com/asciimoo) and the [wuzz](https://github.com/asciimoo/wuzz) project for all the help and suggestions.  

