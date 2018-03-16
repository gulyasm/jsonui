package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type TreePosition []string

func (t TreePosition) Empty() bool {
	return len(t) == 0
}
func (t TreePosition) Shift() TreePosition {
	newLength := len(t) - 1
	newPosition := make([]string, newLength, newLength)
	for i := 0; i < newLength; i += 1 {
		newPosition[i] = t[i+1]
	}
	return newPosition
}

type Query struct {
	q string
}

type TreeNode interface {
	String(int, int) string
	Draw(io.Writer, int, int) error
	Filter(query Query) bool
	Find(TreePosition) TreeNode
	ToggleExpanded()
}

type BaseTreeNode struct {
	isExpanded bool
}

func (n *BaseTreeNode) ToggleExpanded() {
	n.isExpanded = !n.isExpanded
}
func (n BaseTreeNode) expIcon() string {
	if n.isExpanded {
		return "[+]"
	} else {
		return "[-]"
	}
}

type ComplexNode struct {
	BaseTreeNode
	data map[string]TreeNode
}

func (n ComplexNode) Find(tp TreePosition) TreeNode {
	if tp.Empty() {
		return &n
	}
	e, ok := n.data[tp[0]]
	newTp := tp.Shift()
	if !ok {
		// This can't happen in theory
		return nil
	}
	if newTp.Empty() {
		return e
	} else {
		return e.Find(newTp)
	}
}

func (n ComplexNode) stringChildren(padding, lvl int) string {
	s := []string{}
	for key, value := range n.data {
		s = append(s, fmt.Sprintf("%s\"%s\": %s", strings.Repeat(" ", padding+lvl*padding), key, value.String(padding, lvl+1)))
	}
	result := strings.Join(s, ",\n")
	return result
}

func (n ComplexNode) String(padding, lvl int) string {
	return fmt.Sprintf("{\n%s\n%s}", n.stringChildren(padding, lvl), strings.Repeat(" ", lvl*padding))
}

func (n ComplexNode) Draw(writer io.Writer, padding, lvl int) error {
	if lvl == 0 {
		fmt.Fprintf(writer, "%s\n", "root")
	}
	if n.isExpanded {
		for key, value := range n.data {
			fmt.Fprintf(writer, "%s%s\n", strings.Repeat(" ", padding+lvl*padding), key)
			value.Draw(writer, padding, lvl+1)
		}
	} else {
		fmt.Fprintf(writer, "%s%s\n", strings.Repeat(" ", padding+lvl*padding), "+ ...")
	}
	return nil

}
func (n ComplexNode) Filter(query Query) bool {
	return true

}

type ListNode struct {
	BaseTreeNode
	data []TreeNode
}

func (n ListNode) Find(tp TreePosition) TreeNode {
	if tp.Empty() {
		return &n
	}
	i, err := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(tp[0], "["), "]"))
	if err != nil {
		return nil
	}
	newTp := tp.Shift()
	if newTp.Empty() {
		return n.data[i]
	} else {
		return n.data[i].Find(newTp)
	}
}

func (n ListNode) stringChildren(padding, lvl int) string {
	s := []string{}
	for _, value := range n.data {
		s = append(s, strings.Repeat(" ", lvl*padding)+value.String(padding, lvl+1))
	}
	result := strings.Join(s, ",\n")
	return result
}

func (n ListNode) String(padding, lvl int) string {
	return fmt.Sprintf("[\n%s\n%s]", n.stringChildren(padding, lvl+1), strings.Repeat(" ", lvl*padding))
}

func (n ListNode) Draw(writer io.Writer, padding, lvl int) error {
	if lvl == 0 {
		fmt.Fprintf(writer, "%s\n", "root (list)")
	}
	if n.isExpanded {
		for i, value := range n.data {
			fmt.Fprintf(writer, "%s[%d]\n", strings.Repeat(" ", padding+lvl*padding), i)
			value.Draw(writer, padding, lvl+1)
		}
	} else {
		fmt.Fprintf(writer, "%s%s\n", strings.Repeat(" ", padding+lvl*padding), "+ ...")
	}
	return nil

}
func (n ListNode) Filter(query Query) bool {
	return true

}

type FloatNode struct {
	BaseTreeNode
	data float64
}

func (n FloatNode) Find(tp TreePosition) TreeNode {
	return nil

}

func (n FloatNode) String(int, int) string {
	return fmt.Sprintf("%g", n.data)
}

func (n FloatNode) Draw(writer io.Writer, padding, lvl int) error {
	return nil

}
func (n FloatNode) Filter(query Query) bool {
	return true

}

type StringNode struct {
	BaseTreeNode
	data string
}

func (n StringNode) Find(tp TreePosition) TreeNode {
	return nil
}

func (n StringNode) String(_, _ int) string {
	return fmt.Sprintf("%q", n.data)
}

func (n StringNode) Draw(writer io.Writer, padding, lvl int) error {
	//fmt.Fprintf(writer, "%s%q\n", strings.Repeat(" ", padding+padding*lvl), n.data)
	return nil

}
func (n StringNode) Filter(query Query) bool {
	return true
}

func NewTree(y interface{}) (TreeNode, error) {
	var tree TreeNode
	switch v := y.(type) {
	case string:
		tree = &StringNode{
			BaseTreeNode{true},
			v,
		}
	case float64:
		tree = &FloatNode{
			BaseTreeNode{true},
			v,
		}
	case map[string]interface{}:
		data := map[string]TreeNode{}
		for key, childInterface := range v {
			childNode, err := NewTree(childInterface)
			if err != nil {
				return nil, err
			}
			data[key] = childNode
		}
		tree = &ComplexNode{
			BaseTreeNode{true},
			data,
		}
	case []interface{}:
		data := []TreeNode{}
		for _, listItemInterface := range v {
			listItem, err := NewTree(listItemInterface)
			if err != nil {
				return nil, err
			}
			data = append(data, listItem)
		}
		tree = &ListNode{
			BaseTreeNode{true},
			data,
		}
	default:
		tree = &StringNode{BaseTreeNode{true}, "TODO"}

	}
	return tree, nil

}

func FromBytes(b []byte) (TreeNode, error) {
	var y interface{}
	err := json.Unmarshal(b, &y)
	if err != nil {
		log.Fatal("failed to marshal raw json: ", err)
	}
	return NewTree(y)
}

func FromReader(r io.Reader) (TreeNode, error) {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}
	return FromBytes(b)
}
