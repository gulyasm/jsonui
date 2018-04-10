package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"

	"github.com/jroimartin/gocui"
)

const (
	treeView = "tree"
	textView = "text"
	pathView = "path"
	helpView = "help"
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

func logFile(s string) error {
	d1 := []byte(s + "\n")
	return ioutil.WriteFile("log.txt", d1, 0644)
}

func (vp viewPosition) getCoordinates(maxX, maxY int) (int, int, int, int) {
	var x0 = vp.x0.getCoordinate(maxX)
	var y0 = vp.y0.getCoordinate(maxY)
	var x1 = vp.x1.getCoordinate(maxX)
	var y1 = vp.y1.getCoordinate(maxY)
	return x0, y0, x1, y1
}

var helpWindowToggle = false

var viewPositions = map[string]viewPosition{
	treeView: {
		position{0.0, 0},
		position{0.0, 0},
		position{0.3, 2},
		position{0.9, 2},
	},
	textView: {
		position{0.3, 0},
		position{0.0, 0},
		position{1.0, 2},
		position{0.9, 2},
	},
	pathView: {
		position{0.0, 0},
		position{0.89, 0},
		position{1.0, 2},
		position{1.0, 2},
	},
}

var tree treeNode

func main() {
	var err error
	tree, err = fromReader(os.Stdin)
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
	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding(treeView, 'k', gocui.ModNone, cursorMovement(-1)); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(treeView, 'j', gocui.ModNone, cursorMovement(1)); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(treeView, gocui.KeyArrowUp, gocui.ModNone, cursorMovement(-1)); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(treeView, gocui.KeyArrowDown, gocui.ModNone, cursorMovement(1)); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(treeView, 'K', gocui.ModNone, cursorMovement(-15)); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(treeView, 'J', gocui.ModNone, cursorMovement(15)); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(treeView, gocui.KeyPgup, gocui.ModNone, cursorMovement(-15)); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(treeView, gocui.KeyPgdn, gocui.ModNone, cursorMovement(15)); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(treeView, 'e', gocui.ModNone, toggleExpand); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(treeView, 'E', gocui.ModNone, expandAll); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(treeView, 'C', gocui.ModNone, collapseAll); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", 'h', gocui.ModNone, toggleHelp); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", '?', gocui.ModNone, toggleHelp); err != nil {
		log.Panicln(err)
	}
	g.SelFgColor = gocui.ColorBlack
	g.SelBgColor = gocui.ColorGreen

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

const helpMessage = `
JSONUI - Help
----------------------------------------------
j/ArrowDown		═ 	Move a line down
k/ArrowUp 		═ 	Move a line up
J/PageDown		═ 	Move 15 line down
K/PageUp 		═ 	Move 15 line up
e				═ 	Toggle expend/collapse node
E				═ 	Expand all nodes
C				═ 	Collapse all nodes
q/ctrl+c		═ 	Exit
h/?				═ 	Toggle help message
`

func layout(g *gocui.Gui) error {
	var views = []string{treeView, textView, pathView}
	maxX, maxY := g.Size()
	for _, view := range views {
		x0, y0, x1, y1 := viewPositions[view].getCoordinates(maxX, maxY)
		if v, err := g.SetView(view, x0, y0, x1, y1); err != nil {
			v.SelFgColor = gocui.ColorBlack
			v.SelBgColor = gocui.ColorGreen

			v.Title = " " + view + " "
			if err != gocui.ErrUnknownView {
				return err

			}
			if v.Name() == treeView {
				v.Highlight = true
				drawTree(g, v, tree)
				// v.Autoscroll = true
			}
			if v.Name() == textView {
				drawJSON(g, v)
			}

		}
	}
	if helpWindowToggle {
		height := strings.Count(helpMessage, "\n") + 1
		width := -1
		for _, line := range strings.Split(helpMessage, "\n") {
			width = int(math.Max(float64(width), float64(len(line)+2)))
		}
		if v, err := g.SetView(helpView, maxX/2-width/2, maxY/2-height/2, maxX/2+width/2, maxY/2+height/2); err != nil {
			if err != gocui.ErrUnknownView {
				return err

			}
			fmt.Fprintln(v, helpMessage)

		}
	} else {
		g.DeleteView(helpView)
	}
	_, err := g.SetCurrentView(treeView)
	if err != nil {
		log.Fatal("failed to set current view: ", err)
	}
	return nil

}
func getPath(g *gocui.Gui, v *gocui.View) string {
	p := findTreePosition(g)
	for i, s := range p {
		transformed := s
		if !strings.HasPrefix(s, "[") && !strings.HasSuffix(s, "]") {
			transformed = fmt.Sprintf("[%q]", s)
		}
		p[i] = transformed
	}
	return strings.Join(p, "")
}

func drawPath(g *gocui.Gui, v *gocui.View) error {
	pv, err := g.View(pathView)
	if err != nil {
		log.Fatal("failed to get pathView", err)
	}
	p := getPath(g, v)
	pv.Clear()
	fmt.Fprintf(pv, p)
	return nil
}
func drawJSON(g *gocui.Gui, v *gocui.View) error {
	dv, err := g.View(textView)
	if err != nil {
		log.Fatal("failed to get textView", err)
	}
	p := findTreePosition(g)
	treeTodraw := tree.find(p)
	if treeTodraw != nil {
		dv.Clear()
		fmt.Fprintf(dv, treeTodraw.String(2, 0))
	}
	return nil
}

func lineBelow(v *gocui.View, d int) bool {
	_, y := v.Cursor()
	line, err := v.Line(y + d)
	return err == nil && line != ""
}

func countIndex(s string) int {
	count := 0
	for _, c := range s {
		if c == ' ' {
			count++
		}
	}
	return count
}

func getLine(s string, y int) string {
	lines := strings.Split(s, "\n")
	return lines[y]
}

var cleanPatterns = []string{
	treeSignUpEnding,
	treeSignDash,
	treeSignUpMiddle,
	treeSignVertical,
	" (+)",
}

func findTreePosition(g *gocui.Gui) treePosition {
	v, err := g.View(treeView)
	if err != nil {
		log.Fatal("failed to get treeview", err)
	}
	path := treePosition{}
	ci := -1
	_, yOffset := v.Origin()
	_, yCurrent := v.Cursor()
	y := yOffset + yCurrent
	s := v.Buffer()
	for cy := y; cy >= 0; cy-- {
		line := getLine(s, cy)
		for _, pattern := range cleanPatterns {
			line = strings.Replace(line, pattern, "", -1)
		}

		if count := countIndex(line); count < ci || ci == -1 {
			path = append(path, strings.TrimSpace(line))
			ci = count
		}
	}
	for i := len(path)/2 - 1; i >= 0; i-- {
		opp := len(path) - 1 - i
		path[i], path[opp] = path[opp], path[i]
	}

	return path[1:]
}

// This is a workaround for not having a Buffer
// function in gocui
func bufferLen(v *gocui.View) int {
	s := v.Buffer()
	return len(strings.Split(s, "\n")) - 1
}

func drawTree(g *gocui.Gui, v *gocui.View, tree treeNode) error {
	tv, err := g.View(treeView)
	if err != nil {
		log.Fatal("failed to get treeView", err)
	}
	tv.Clear()
	tree.draw(tv, 2, 0)
	maxY := bufferLen(tv)
	cx, cy := tv.Cursor()
	lastLine := maxY - 2
	if cy > lastLine {
		tv.SetCursor(cx, lastLine)
		tv.SetOrigin(0, 0)
	}

	return nil
}

func expandAll(g *gocui.Gui, v *gocui.View) error {
	tree.expandAll()
	return drawTree(g, v, tree)
}

func collapseAll(g *gocui.Gui, v *gocui.View) error {
	tree.collapseAll()
	return drawTree(g, v, tree)
}

func toggleExpand(g *gocui.Gui, v *gocui.View) error {
	p := findTreePosition(g)
	subTree := tree.find(p)
	subTree.toggleExpanded()
	return drawTree(g, v, tree)
}

func cursorMovement(d int) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		dir := 1
		if d < 0 {
			dir = -1
		}
		distance := int(math.Abs(float64(d)))
		for ; distance > 0; distance-- {
			if lineBelow(v, distance*dir) {
				v.MoveCursor(0, distance*dir, false)
				drawJSON(g, v)
				drawPath(g, v)
				return nil
			}
		}
		return nil
	}
}
func toggleHelp(g *gocui.Gui, v *gocui.View) error {
	helpWindowToggle = !helpWindowToggle
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
