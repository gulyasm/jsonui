# JSONUI
[![](https://travis-ci.org/gulyasm/jsonui.svg?branch=master)](https://travis-ci.org/gulyasm/jsonui) [![](https://goreportcard.com/badge/github.com/gulyasm/jsonui)](https://goreportcard.com/report/github.com/gulyasm/jsonui)

`jsonui` is an interactive JSON explorer in your command line. You can pipe any JSON into `jsonui` and explore it, copy the path for each element.

![](img/jsonui.gif)

## Import into your application
jsonui can now be imported into your application --
```
...
jsonBytes,err := ioutil.ReadFile("jsonfile.json")
errorhandle(err)
jsonpath:=jsonui.Interactive(jsonBytes)
fmt.Println(jsonpath)
...
```
When it runs, scroll down to the section, press x to execute. The output will look like this:
```
] $ cat test.json | jsonui
JSON Path: address.gateways
```
And then the user will be able to extract json using another utility called 'jj'
```
]$ cat test.json | jj address.gateways
["Sopron", "Vienna", "Budapest"]
```

You could also use the jsonpath variable above for other packages like gjson or sjson to get / set values.
## Install
This Version:
`git clone https://github.com/rmasci/jsonui.git`
cd jsonui/jsonui-cmd
make

Original:
`go get -u github.com/gulyasm/jsonui`

## Binary Releases
This Version:
[Binary releases are availabe](https://github.com/rmasci/jsonui/releases)

Original:
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
### `x`
Execute. Will give you the JSON Path to use with jj command.
[https://github.com/tidwall/jj](https://github.com/tidwall/jj)
#### `q/Ctrl+C`
Quit jsonui



## Acknowledgments
Special thanks for [asciimoo](https://github.com/asciimoo) and the [wuzz](https://github.com/asciimoo/wuzz) project for all the help and suggestions.  

