# JSONUI
[![](https://travis-ci.org/gulyasm/jsonui.svg?branch=master)](https://travis-ci.org/gulyasm/jsonui) [![](https://goreportcard.com/badge/github.com/gulyasm/jsonui)](https://goreportcard.com/report/github.com/gulyasm/jsonui)

`jsonui` is an interactive JSON explorer in your command line. You can pipe any JSON into `jsonui` and explore it, copy the path for each element.

![](img/screenshot1.png)
![](img/screenshot2.png)

## Install
`go get github.com/gulyasm.jsonui`

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

#### `J`
Move down 15 lines

#### `K`
Move up 15 lines

#### `e`
Toggle node (expend or collapse)
