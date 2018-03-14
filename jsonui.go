package main

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/jroimartin/gocui"
)

const (
	TREE_VIEW = "view/tree"
	TEXT_VIEW = "view/text"
)

type position struct {
	prc    float32
	margin int
}

func (p position) getCoordinate(max int) int {
	// value = prc * MAX + abs
	return int(p.prc*float32(max)) - p.margin
}

type viewPosition struct {
	x0, y0, x1, y1 position
}

func (vp viewPosition) getCoordinates(maxX, maxY int) (int, int, int, int) {
	var x0 = vp.x0.getCoordinate(maxX)
	var y0 = vp.y0.getCoordinate(maxY)
	var x1 = vp.x1.getCoordinate(maxX)
	var y1 = vp.y1.getCoordinate(maxY)
	return x0, y0, x1, y1
}

var VIEW_POSITIONS = map[string]viewPosition{
	TREE_VIEW: {
		position{0.0, 0},
		position{0.0, 0},
		position{0.3, 2},
		position{1.0, 2},
	},
	TEXT_VIEW: {
		position{0.3, 0},
		position{0.0, 0},
		position{1.0, 2},
		position{1.0, 2},
	},
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)

	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)

	}
}

type TreeNode struct {
	text     string
	children []*TreeNode
}

var mytree = TreeNode{
	"root",
	[]*TreeNode{
		&TreeNode{"hello1", []*TreeNode{
			&TreeNode{"mam", nil},
			&TreeNode{"mam", nil},
			&TreeNode{"papapa", nil},
			&TreeNode{"mam", nil},
			&TreeNode{"papapa", nil},
			&TreeNode{"mam", nil},
			&TreeNode{"papapa", nil},
			&TreeNode{"mam", nil},
			&TreeNode{"papapa", nil},
			&TreeNode{"papapa", nil},
		}},
		&TreeNode{"hello2", nil},
	},
}

func (node *TreeNode) Draw(writer io.Writer, lvl int) error {
	fmt.Fprintln(writer, strings.Repeat("  ", lvl)+" "+node.text)
	for _, child := range node.children {
		err := child.Draw(writer, lvl+1)
		if err != nil {
			return err
		}
	}
	return nil
}

func layout(g *gocui.Gui) error {
	var views = []string{TREE_VIEW, TEXT_VIEW}
	maxX, maxY := g.Size()
	for _, view := range views {
		x0, y0, x1, y1 := VIEW_POSITIONS[view].getCoordinates(maxX, maxY)
		if v, err := g.SetView(view, x0, y0, x1, y1); err != nil {
			if err != gocui.ErrUnknownView {
				return err

			}
			mytree.Draw(v, 0)

		}
	}
	return nil

}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit

}
