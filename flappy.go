package main

import (
	"log"
	"time"

	"github.com/jroimartin/gocui"
)

// point includes coordinates of a point
type point struct {
	x int
	y int
}

// BirdPosition flappy bird position on view
var BirdPosition point

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(launchGame)

	if err := g.SetKeybinding("", gocui.KeySpace, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, birdUp); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}

func launchGame(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("flappy", 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.FgColor = gocui.ColorYellow
		//		ticker := time.NewTicker(500 * time.Millisecond)
		BirdPosition = point{x: 10, y: 10}

		go func() {
			for true {
				birdDown(g, v)
				refreshBoard(g, v)
				time.Sleep(200 * time.Millisecond)
			}
		}()

	}

	return nil
}

func birdUp(g *gocui.Gui, v *gocui.View) error {
	BirdPosition.y = BirdPosition.y - 4
	return nil
}

func birdDown(g *gocui.Gui, v *gocui.View) error {
	BirdPosition.y++
	return nil
}

func refreshBoard(g *gocui.Gui, v *gocui.View) error {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("flappy")
		if err != nil {

		}
		v.Clear()
		moveBird(BirdPosition, v)
		drawBorders(v)
		return nil
	})
	return nil
}

func moveBird(p point, v *gocui.View) {
	v.SetCursor(p.x, p.y)
	v.EditWrite('@')

}

func drawBorders(v *gocui.View) error {
	vMaxX, vMaxY := v.Size()

	// draw top and bottom edge borders
	drawPipe(point{x: 0, y: 0}, point{x: vMaxX - 1, y: 0}, v)
	drawPipe(point{x: 0, y: vMaxY - 1}, point{x: vMaxX - 1, y: vMaxY - 1}, v)

	return nil
}

func drawPipe(p1 point, p2 point, v *gocui.View) {

	for y := p1.y; y <= p2.y; y++ {
		for x := p1.x; x <= p2.x; x++ {
			v.SetCursor(x, y)
			v.EditWrite('#')
		}
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
