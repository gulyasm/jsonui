package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Query struct {
	q string
}

type TreeNode interface {
	String(int, int) string
	Draw(io.Writer, int, int) error
	Filter(query Query) bool
}

type BaseTreeNode struct {
	isExpanded bool
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

func (n ComplexNode) stringChildren(padding, lvl int) string {
	s := []string{}
	for key, value := range n.data {
		s = append(s, fmt.Sprintf("%s\"%s\": %s", strings.Repeat(" ", lvl*padding), key, value.String(padding, lvl+1)))
	}
	result := strings.Join(s, ",\n")
	return result
}

func (n ComplexNode) String(padding, lvl int) string {
	return fmt.Sprintf("{\n%s\n%s}", n.stringChildren(padding, lvl), strings.Repeat(" ", lvl*padding))
}

func (n ComplexNode) Draw(writer io.Writer, padding, lvl int) error {
	for key, value := range n.data {
		fmt.Fprintf(writer, "%s%s\n", strings.Repeat(" ", lvl*padding), key)
		value.Draw(writer, padding, lvl+1)
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
	for i, value := range n.data {
		fmt.Fprintf(writer, "%s[%d]", strings.Repeat(" ", lvl*padding), i)
		value.Draw(writer, padding, lvl+1)
		fmt.Fprintf(writer, "\n")
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

func main() {
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	var y map[string]interface{}
	json.Unmarshal(bytes, &y)
	tree, err := NewTree(y)
	if err != nil {
		log.Fatal(err)
	}
	tree.Draw(os.Stdout, 4, 0)
}
