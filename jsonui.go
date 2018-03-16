package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/jroimartin/gocui"
)

const (
	TREE_VIEW = "tree"
	TEXT_VIEW = "text"
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
		position{0.2, 2},
		position{1.0, 2},
	},
	TEXT_VIEW: {
		position{0.2, 0},
		position{0.0, 0},
		position{1.0, 2},
		position{1.0, 2},
	},
}

var tree TreeNode

func main() {
	var err error
	tree, err = FromReader(os.Stdin)
	if err != nil {
		log.Panicln(err)
	}
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding(TREE_VIEW, 'k', gocui.ModNone, cursorUp); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(TREE_VIEW, 'j', gocui.ModNone, cursorDown); err != nil {
		log.Panicln(err)
	}

	g.SelFgColor = gocui.ColorBlack
	g.SelBgColor = gocui.ColorGreen

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	var views = []string{TREE_VIEW, TEXT_VIEW}
	maxX, maxY := g.Size()
	for _, view := range views {
		x0, y0, x1, y1 := VIEW_POSITIONS[view].getCoordinates(maxX, maxY)
		if v, err := g.SetView(view, x0, y0, x1, y1); err != nil {
			v.SelFgColor = gocui.ColorBlack
			v.SelBgColor = gocui.ColorGreen

			v.Title = " " + view + " "
			if err != gocui.ErrUnknownView {
				return err

			}
			if v.Name() == TREE_VIEW {
				v.Highlight = true
				tree.Draw(v, 2, 0)
			}
		}
	}
	_, err := g.SetCurrentView(TREE_VIEW)
	if err != nil {
		log.Fatal(err)
	}
	return nil

}
func drawText(g *gocui.Gui, v *gocui.View) error {
	textView, err := g.View(TEXT_VIEW)
	if err != nil {
		log.Fatal(err)
	}
	textView.Clear()
	fmt.Fprintf(textView, tree.String(2, 0))
	return nil
}

func lineBelow(v *gocui.View) bool {
	_, y := v.Cursor()
	line, err := v.Line(y + 1)
	return err == nil && line != ""
}

func CountIndent(s string) int {
	count := 0
	for _, c := range s {
		if c == ' ' {
			count += 1
		}
	}
	return count
}

func FindTreePosition(v *gocui.View, dv io.Writer) TreePosition {
	path := TreePosition{}
	ci := -1
	for _, cy := v.Cursor(); cy >= 0; cy -= 1 {
		line, err := v.Line(cy)
		if err != nil {
			log.Fatal(err)
		}
		if count := CountIndent(line); count < ci || ci == -1 {
			path = append(path, strings.TrimSpace(line))
			ci = count
		}
	}
	for i := len(path)/2 - 1; i >= 0; i-- {
		opp := len(path) - 1 - i
		path[i], path[opp] = path[opp], path[i]

	}

	for _, p := range path {
		fmt.Fprintln(dv, p)
	}
	return path
}

func debugView(g *gocui.Gui) *gocui.View {
	textView, err := g.View(TEXT_VIEW)
	if err != nil {
		log.Fatal(err)
	}
	textView.Clear()
	return textView
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	dv := debugView(g)
	if lineBelow(v) {
		v.MoveCursor(0, 1, false)
		FindTreePosition(v, dv)
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	dv := debugView(g)
	if lineBelow(v) {
		v.MoveCursor(0, -1, false)
		FindTreePosition(v, dv)
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
