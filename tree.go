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
	// Draw(io.Writer, int, int) error
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
		s = append(s, fmt.Sprintf("%s\"%s\": %s", strings.Repeat(" ", padding+lvl*padding), key, value.String(padding, lvl+1)))
	}
	result := strings.Join(s, ",\n")
	return result
}

func (n ComplexNode) String(padding, lvl int) string {
	return fmt.Sprintf("{\n%s\n%s}", n.stringChildren(padding, lvl), strings.Repeat(" ", lvl*padding))
}

func (n ComplexNode) Draw(io.Writer, int, int) error {
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
		s = append(s, strings.Repeat(" ", padding+lvl*padding)+value.String(padding, lvl+1))
	}
	result := strings.Join(s, ",\n")
	return result
}

func (n ListNode) String(padding, lvl int) string {
	return fmt.Sprintf("[\n%s\n%s]", n.stringChildren(padding, lvl+1), strings.Repeat(" ", lvl*padding))
}

func (n ListNode) Draw(io.Writer, int, int) error {
	return nil

}
func (n ListNode) Filter(query Query) bool {
	return true

}

type FloatNode struct {
	BaseTreeNode
	data float64
}

func (n FloatNode) String(_, lvl int) string {
	return fmt.Sprintf("%g", n.data)
}

func (n FloatNode) Draw(io.Writer, int, int) error {
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

func (n StringNode) Draw(io.Writer, int, int) error {
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

/*

func (node *TreeNode) Draw(writer io.Writer, lvl, padding int) error {
	str := fmt.Sprintf("%-"+strconv.Itoa(padding)+"s",
		strings.Repeat("  ", lvl)+" "+node.String())
	fmt.Fprintln(writer, str)
	if node.isExpanded {
		for _, child := range node.children {
			err := child.Draw(writer, lvl+1, padding)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func newRecursiveTree(key string, y interface{}) (*TreeNode, error) {
	var tree *TreeNode
	switch v := y.(type) {
	case int:
		tree = &TreeNode{key, strconv.Itoa(v), TypeInt, nil, true}
	case float64:
		tree = &TreeNode{key, strconv.FormatFloat(v, 'f', 6, 64), TypeFloat, nil, true}
	case string:
		tree = &TreeNode{key, v, TypeString, nil, true}
	case map[string]interface{}:
		children := []*TreeNode{}
		for key, child := range v {
			childNode, err := newRecursiveTree(key, child)
			if err != nil {
				return nil, err
			}
			children = append(children, childNode)
		}
		tree = &TreeNode{key, key, TypeComplex, children, true}
	case []interface{}:
		children := []*TreeNode{}
		for i, child := range v {
			childNode, err := newRecursiveTree(strconv.Itoa(i), child)
			if err != nil {
				return nil, err
			}
			children = append(children, childNode)
		}
		tree = &TreeNode{key, "", TypeList, children, true}
	default:
		tree = &TreeNode{key, "TODO", -1, nil, true}

	}
	return tree, nil
}
*/
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
	str := tree.String(2, 0)
	fmt.Println(str)
}
