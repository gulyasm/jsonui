package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/jroimartin/gocui"
)

const (
	TREE_VIEW = "tree"
	TEXT_VIEW = "text"
	PATH_VIEW = "path"
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

func LogFile(s string) error {
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

var VIEW_POSITIONS = map[string]viewPosition{
	TREE_VIEW: {
		position{0.0, 0},
		position{0.0, 0},
		position{0.2, 2},
		position{0.9, 2},
	},
	TEXT_VIEW: {
		position{0.2, 0},
		position{0.0, 0},
		position{1.0, 2},
		position{0.9, 2},
	},
	PATH_VIEW: {
		position{0.0, 0},
		position{0.89, 0},
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

	if err := g.SetKeybinding(TREE_VIEW, 'k', gocui.ModNone, cursorMovement(-1)); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(TREE_VIEW, 'j', gocui.ModNone, cursorMovement(1)); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(TREE_VIEW, 'K', gocui.ModNone, cursorMovement(-15)); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(TREE_VIEW, 'J', gocui.ModNone, cursorMovement(15)); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(TREE_VIEW, 'e', gocui.ModNone, toggleExpand); err != nil {
		log.Panicln(err)
	}
	g.SelFgColor = gocui.ColorBlack
	g.SelBgColor = gocui.ColorGreen

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	var views = []string{TREE_VIEW, TEXT_VIEW, PATH_VIEW}
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
				// v.Autoscroll = true
			}
			if v.Name() == TEXT_VIEW {
				drawJson(g, v)
			}

		}
	}
	_, err := g.SetCurrentView(TREE_VIEW)
	if err != nil {
		log.Fatal("failed to set current view: ", err)
	}
	return nil

}
func getPath(g *gocui.Gui, v *gocui.View) string {
	tv, err := g.View(TREE_VIEW)
	if err != nil {
		log.Fatal("failed to get TREE_VIEW", err)
	}
	p := FindTreePosition(tv)
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
	pv, err := g.View(PATH_VIEW)
	if err != nil {
		log.Fatal("failed to get PATH_VIEW", err)
	}
	p := getPath(g, v)
	pv.Clear()
	fmt.Fprintf(pv, p)
	return nil
}
func drawJson(g *gocui.Gui, v *gocui.View) error {
	dv, err := g.View(TEXT_VIEW)
	if err != nil {
		log.Fatal("failed to get TEXT_VIEW", err)
	}
	tv, err := g.View(TREE_VIEW)
	if err != nil {
		log.Fatal("failed to get TREE_VIEW", err)
	}
	p := FindTreePosition(tv)
	LogFile(strings.Join(p, " / "))
	treeToDraw := tree.Find(p)
	if treeToDraw != nil {
		dv.Clear()
		fmt.Fprintf(dv, treeToDraw.String(2, 0))
	}
	return nil
}

func lineBelow(v *gocui.View, d int) bool {
	_, y := v.Cursor()
	line, err := v.Line(y + d)
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

func GetLine(s string, y int) string {
	lines := strings.Split(s, "\n")
	return lines[y]
}

func FindTreePosition(v *gocui.View) TreePosition {
	path := TreePosition{}
	ci := -1
	lg := ""
	_, yOffset := v.Origin()
	_, yCurrent := v.Cursor()
	y := yOffset + yCurrent
	s := v.Buffer()
	for cy := y; cy >= 0; cy -= 1 {
		line := GetLine(s, cy)

		if strings.TrimSpace(line) == "+ ..." {
			continue
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

	return path[1:]
}

func toggleExpand(g *gocui.Gui, v *gocui.View) error {
	tv, err := g.View(TREE_VIEW)
	if err != nil {
		log.Fatal("failed to get TREE_VIEW", err)
	}
	p := FindTreePosition(tv)
	subTree := tree.Find(p)
	subTree.ToggleExpanded()
	tv.Clear()
	tree.Draw(tv, 2, 0)
	return nil
}
func cursorMovement(d int) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if lineBelow(v, d) {
			v.MoveCursor(0, d, false)
			drawJson(g, v)
			drawPath(g, v)
		}
		return nil
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
