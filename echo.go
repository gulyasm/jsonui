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

type TreeNode struct {
	text     string
	children []*TreeNode
}

func (node *TreeNode) Draw(writer io.Writer, lvl, padding int) error {
	str := fmt.Sprintf("%-"+strconv.Itoa(padding)+"s", strings.Repeat("  ", lvl)+" "+node.text)
	fmt.Fprintln(writer, str)
	for _, child := range node.children {
		err := child.Draw(writer, lvl+1, padding)
		if err != nil {
			return err
		}
	}
	return nil
}

func newRecursiveTree(y interface{}) (*TreeNode, error) {
	var tree *TreeNode
	switch v := y.(type) {
	case int:
		tree = &TreeNode{strconv.Itoa(v), nil}
	case float64:
		tree = &TreeNode{strconv.FormatFloat(v, 'f', 6, 64), nil}
	case string:
		tree = &TreeNode{v, nil}
	default:
		tree = &TreeNode{"TODO", nil}
		// i isn't one of the types above

	}
	return tree, nil
}

func NewTree(y map[string]interface{}) (*TreeNode, error) {
	//newRecursiveTree(tree)
	children := []*TreeNode{}
	for _, child := range y {
		childNode, err := newRecursiveTree(child)
		if err != nil {
			return nil, err
		}
		children = append(children, childNode)
	}
	root := &TreeNode{"root", children}
	return root, nil
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
	tree.Draw(os.Stdout, 0, 12)
	return
	log.Println(y)
	str, err := json.MarshalIndent(y, "", " ")
	if err != nil {
		fmt.Println("Error encoding JSON")
		return

	}

	fmt.Println(string(str))
}
