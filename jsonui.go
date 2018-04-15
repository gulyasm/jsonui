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

const VERSION = "1.0.1"

const (
	searchView = "search"
	treeView   = "tree"
	textView   = "text"
	pathView   = "path"
	helpView   = "help"
)

func logFile(s string) error {
	d1 := []byte(s + "\n")
	return ioutil.WriteFile("log.txt", d1, 0644)
}

var helpWindowToggle = false

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
	if err := g.SetKeybinding("", gocui.KeyCtrlT, gocui.ModNone, focusView(treeView)); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlS, gocui.ModNone, focusView(searchView)); err != nil {
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

type DefaultEditor struct {
	g *gocui.Gui
}

func (d *DefaultEditor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
		search(v.Buffer(), d.g, v)
	case key == gocui.KeySpace:
		return
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
		search(v.Buffer(), d.g, v)
	default:
		return
	}

}

func ratio(ratio float32, max int) int {
	return int(float32(max) * ratio)
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(searchView, 0, 0, ratio(0.6, maxX)-2, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.SelFgColor = gocui.ColorBlack
		v.SelBgColor = gocui.ColorGreen
		v.Editor = &DefaultEditor{g}
		v.Editable = true

		v.Title = " " + searchView + " "

	}
	if v, err := g.SetView(treeView, 0, 3, ratio(0.6, maxX)-2, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " " + treeView + " "
		v.Highlight = true
		drawTree(g, v, tree)
	}
	if v, err := g.SetView(pathView, 0, maxY-3, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " " + pathView + " "
	}
	if v, err := g.SetView(textView, ratio(0.6, maxX), 0, maxX-1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " " + textView + " "
		drawJSON(g, v)
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
	g.Highlight = true
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

func focusView(view string) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		g.SetCurrentView(view)
		return nil
	}
}

func collapseAll(g *gocui.Gui, v *gocui.View) error {
	tree.collapseAll()
	return drawTree(g, v, tree)
}

func search(query string, g *gocui.Gui, v *gocui.View) error {
	query = strings.TrimSuffix(query, "\n")
	filteredTree, _ := tree.search(query)
	if filteredTree != nil {
		drawTree(g, v, filteredTree)
	}
	return nil
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
