package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	treeSignDash     = "─"
	treeSignVertical = "│"
	treeSignUpMiddle = "├"
	treeSignUpEnding = "└"
)

type treePosition []string

func (t treePosition) empty() bool {
	return len(t) == 0
}
func (t treePosition) shift() treePosition {
	newLength := len(t) - 1
	newPosition := make([]string, newLength, newLength)
	for i := 0; i < newLength; i++ {
		newPosition[i] = t[i+1]
	}
	return newPosition
}

type query struct {
	q string
}

type treeNode interface {
	String(int, int) string
	draw(io.Writer, int, int) error
	filter(query query) bool
	find(treePosition) treeNode
	search(query string) (treeNode, error)
	isCollapsable() bool
	toggleExpanded()
	collapseAll()
	expandAll()
	isExpanded() bool
}

type baseTreeNode struct {
	expanded bool
}

func (n *baseTreeNode) isExpanded() bool {
	return n.expanded
}

func (n *baseTreeNode) toggleExpanded() {
	n.expanded = !n.expanded
}

func (n baseTreeNode) expIcon() string {
	if n.expanded {
		return "[+]"
	}
	return "[-]"
}

type complexNode struct {
	baseTreeNode
	data map[string]treeNode
}

func (n *complexNode) collapseAll() {
	n.expanded = false
	for _, v := range n.data {
		v.collapseAll()
	}
}
func (n *complexNode) expandAll() {
	n.expanded = true
	for _, v := range n.data {
		v.expandAll()
	}
}
func (n complexNode) isCollapsable() bool {
	return true
}
func (n complexNode) search(query string) (treeNode, error) {
	filteredNode := &complexNode{
		baseTreeNode{true},
		map[string]treeNode{},
	}
	for key, value := range n.data {
		if key == query {
			filteredNode.data[key] = value
		}
	}
	return filteredNode, nil
}
func (n complexNode) find(tp treePosition) treeNode {
	if tp.empty() {
		return &n
	}
	e, ok := n.data[tp[0]]
	newTp := tp.shift()
	if !ok {
		// This can't happen in theory
		return nil
	}
	if newTp.empty() {
		return e
	}
	return e.find(newTp)
}

func (n complexNode) stringChildren(padding, lvl int) string {
	s := []string{}
	keys := []string{}
	for key := range n.data {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		value, _ := n.data[key]
		s = append(s, fmt.Sprintf("%s\"%s\": %s", strings.Repeat(" ", padding+lvl*padding), key, value.String(padding, lvl+1)))
	}
	result := strings.Join(s, ",\n")
	return result
}

func (n complexNode) String(padding, lvl int) string {
	return fmt.Sprintf("{\n%s\n%s}", n.stringChildren(padding, lvl), strings.Repeat(" ", lvl*padding))
}

func (n complexNode) draw(writer io.Writer, padding, lvl int) error {
	if lvl == 0 {
		fmt.Fprintf(writer, "%s\n", "root")
	}
	keys := []string{}
	for key := range n.data {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	length := len(keys)
	for i, key := range keys {
		value, _ := n.data[key]
		var char string
		if i < length-1 {
			char = treeSignUpMiddle
		} else {
			char = treeSignUpEnding
		}
		char += treeSignDash
		expendedCharacter := ""
		if value.isCollapsable() && !value.isExpanded() {
			expendedCharacter += " (+)"
		}
		fmt.Fprintf(writer,
			"%s%s %s%s\n",
			strings.Repeat("│  ", lvl),
			char,
			key,
			expendedCharacter,
		)
		if value.isExpanded() {
			value.draw(writer, padding, lvl+1)
		}
	}
	return nil

}
func (n complexNode) filter(query query) bool {
	return true

}

type listNode struct {
	baseTreeNode
	data []treeNode
}

func (n *listNode) collapseAll() {
	n.expanded = false
	for _, v := range n.data {
		v.collapseAll()
	}
}
func (n *listNode) expandAll() {
	n.expanded = true
	for _, v := range n.data {
		v.expandAll()
	}
}
func (n listNode) isCollapsable() bool {
	return true
}

func (n listNode) search(query string) (treeNode, error) {
	return nil, nil

}
func (n listNode) find(tp treePosition) treeNode {
	if tp.empty() {
		return &n
	}
	i, err := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(tp[0], "["), "]"))
	if err != nil {
		return nil
	}
	newTp := tp.shift()
	if newTp.empty() {
		return n.data[i]
	}
	return n.data[i].find(newTp)
}

func (n listNode) stringChildren(padding, lvl int) string {
	s := []string{}
	for _, value := range n.data {
		s = append(s, strings.Repeat(" ", lvl*padding)+value.String(padding, lvl+1))
	}
	result := strings.Join(s, ",\n")
	return result
}

func (n listNode) String(padding, lvl int) string {
	return fmt.Sprintf("[\n%s\n%s]", n.stringChildren(padding, lvl+1), strings.Repeat(" ", lvl*padding))
}

func (n listNode) draw(writer io.Writer, padding, lvl int) error {
	if lvl == 0 {
		fmt.Fprintf(writer, "%s\n", "root")
	}
	length := len(n.data)
	for i, value := range n.data {
		var char string
		if i < length-1 {
			char = treeSignUpMiddle
		} else {
			char = treeSignUpEnding
		}
		char += treeSignDash
		expendedCharacter := ""

		if value.isCollapsable() && !value.isExpanded() {
			expendedCharacter += " (+)"
		}

		fmt.Fprintf(writer,
			"%s%s [%d]%s\n",
			strings.Repeat("│  ", lvl),
			char,
			i,
			expendedCharacter,
		)
		if value.isExpanded() {
			value.draw(writer, padding, lvl+1)
		}
	}
	return nil

}
func (n listNode) filter(query query) bool {
	return true

}

type floatNode struct {
	baseTreeNode
	data float64
}

func (n *floatNode) collapseAll() {
}
func (n *floatNode) expandAll() {
}
func (n floatNode) isCollapsable() bool {
	return false
}
func (n floatNode) search(query string) (treeNode, error) {
	return nil, nil
}
func (n floatNode) find(tp treePosition) treeNode {
	return nil

}

func (n floatNode) String(int, int) string {
	return fmt.Sprintf("%g", n.data)
}

func (n floatNode) draw(writer io.Writer, padding, lvl int) error {
	return nil

}
func (n floatNode) filter(query query) bool {
	return true
}

type stringNode struct {
	baseTreeNode
	data string
}

func (n *stringNode) collapseAll() {
}
func (n *stringNode) expandAll() {
}

func (n stringNode) isCollapsable() bool {
	return false
}

func (n stringNode) find(tp treePosition) treeNode {
	return nil
}

func (n stringNode) String(_, _ int) string {
	return fmt.Sprintf("%q", n.data)
}

func (n stringNode) search(query string) (treeNode, error) {
	return nil, nil

}
func (n stringNode) draw(writer io.Writer, padding, lvl int) error {
	//fmt.Fprintf(writer, "%s%q\n", strings.Repeat(" ", padding+padding*lvl), n.data)
	return nil

}
func (n stringNode) filter(query query) bool {
	return true
}

type boolNode struct {
	baseTreeNode
	data bool
}

func (n *boolNode) collapseAll() {
}
func (n *boolNode) expandAll() {
}

func (n boolNode) isCollapsable() bool {
	return false
}

func (n boolNode) find(tp treePosition) treeNode {
	return nil
}

func (n boolNode) String(_, _ int) string {
	return fmt.Sprintf("%t", n.data)
}

func (n boolNode) search(query string) (treeNode, error) {
	return nil, nil

}
func (n boolNode) draw(writer io.Writer, padding, lvl int) error {
	return nil

}
func (n boolNode) filter(query query) bool {
	return true
}

type nilNode struct {
	baseTreeNode
}

func (n *nilNode) collapseAll() {
}
func (n *nilNode) expandAll() {
}

func (n nilNode) isCollapsable() bool {
	return false
}

func (n nilNode) find(tp treePosition) treeNode {
	return nil
}

func (n nilNode) String(_, _ int) string {
	return "null"
}

func (n nilNode) search(query string) (treeNode, error) {
	return nil, nil

}
func (n nilNode) draw(writer io.Writer, padding, lvl int) error {
	return nil

}
func (n nilNode) filter(query query) bool {
	return true
}
func newTree(y interface{}) (treeNode, error) {
	var tree treeNode
	switch v := y.(type) {
	case bool:
		tree = &boolNode{
			baseTreeNode{true},
			v,
		}
	case string:
		tree = &stringNode{
			baseTreeNode{true},
			v,
		}

	case nil:
		tree = &nilNode{baseTreeNode{true}}

	case float64:
		tree = &floatNode{
			baseTreeNode{true},
			v,
		}
	case map[string]interface{}:
		data := map[string]treeNode{}
		for key, childInterface := range v {
			childNode, err := newTree(childInterface)
			if err != nil {
				return nil, err
			}
			data[key] = childNode
		}
		tree = &complexNode{
			baseTreeNode{true},
			data,
		}
	case []interface{}:
		data := []treeNode{}
		for _, listItemInterface := range v {
			listItem, err := newTree(listItemInterface)
			if err != nil {
				return nil, err
			}
			data = append(data, listItem)
		}
		tree = &listNode{
			baseTreeNode{true},
			data,
		}
	default:
		tree = &stringNode{baseTreeNode{true}, "TODO"}

	}
	return tree, nil

}

func fromBytes(b []byte) (treeNode, error) {
	var y interface{}
	err := json.Unmarshal(b, &y)
	if err != nil {
		log.Fatal("failed to marshal raw json: ", err)
	}
	return newTree(y)
}

func fromReader(r io.Reader) (treeNode, error) {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}
	return fromBytes(b)
}
