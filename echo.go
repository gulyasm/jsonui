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

type nodeType int

const (
	TypeComplex = nodeType(0)
	TypeList    = iota
	TypeString  = iota
	TypeInt     = iota
	TypeFloat   = iota
)

type TreeNode struct {
	key        string
	value      string
	nodeType   nodeType
	children   []*TreeNode
	isExpanded bool
}

func (node TreeNode) String() string {
	switch node.nodeType {
	case TypeInt:
		return fmt.Sprintf("%s -> %s", node.key, node.value)
	case TypeFloat:
		return fmt.Sprintf("%s -> %s", node.key, node.value)
	case TypeList:
		return fmt.Sprintf("%s [ ", node.key)
	case TypeComplex:
		return fmt.Sprintf("%s {", node.key)
	case TypeString:
		return fmt.Sprintf("%s -> %s", node.key, node.value)
	default:
		log.Fatal("Unknown type")
	}
	return ""

}

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
		for _, child := range v {
			childNode, err := newRecursiveTree("", child)
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

func main() {
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	var y map[string]interface{}
	json.Unmarshal(bytes, &y)
	tree, err := newRecursiveTree("", y)
	if err != nil {
		log.Fatal(err)
	}
	tree.Draw(os.Stdout, 0, 12)
	return
	str, err := json.MarshalIndent(y, "", " ")
	if err != nil {
		fmt.Println("Error encoding JSON")
		return

	}

	fmt.Println(string(str))
}
