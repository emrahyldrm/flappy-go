package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/jroimartin/gocui"
)

type point struct {
	x int
	y int
}

type pipeLocation struct {
	tlFirst  point
	brFirst  point
	tlSecond point
	brSecond point
}

var birdPosition point
var vMaxX int
var vMaxY int
var pipeLocations [2]pipeLocation
var rng *rand.Rand
var pipePieceUnit int

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

// main function of game
func launchGame(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("flappy", 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		birdPosition = point{x: 10, y: 10}
		v.FgColor = gocui.ColorYellow
		vMaxX, vMaxY = v.Size()
		rs := rand.NewSource(54)
		rng = rand.New(rs)
		pipePieceUnit = (vMaxY - 2) / 9
		pipeLocations[0] = calculateNewPipePosition()

		ticker := time.NewTicker(200 * time.Millisecond)
		go func() {
			for range ticker.C {
				birdDown(g, v)
				refreshBoard(g, v)
			}
		}()

	}

	return nil
}

// a step up the position of bird
func birdUp(g *gocui.Gui, v *gocui.View) error {
	birdPosition.y = birdPosition.y - 4
	return nil
}

// a step down the position of bird
func birdDown(g *gocui.Gui, v *gocui.View) error {
	birdPosition.y++
	return nil
}

// update canvas
func refreshBoard(g *gocui.Gui, v *gocui.View) error {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("flappy")
		if err != nil {

		}
		v.Clear()
		drawBird(birdPosition, v)
		drawBorders(v)
		drawPipe(v)
		return nil
	})
	return nil
}

// put the bird at desired point
func drawBird(p point, v *gocui.View) {
	v.SetCursor(p.x, p.y)
	v.EditWrite('@')

}

func calculateNewPipePosition() pipeLocation {
	var pipeLoc pipeLocation
	firstPieceVolume := rng.Intn(8)

	pipeLoc.tlFirst.x = vMaxX - 18
	pipeLoc.tlFirst.y = 1
	pipeLoc.brFirst.x = vMaxX - 15
	pipeLoc.brFirst.y = pipeLoc.tlFirst.y + (firstPieceVolume * pipePieceUnit)

	pipeLoc.tlSecond.x = vMaxX - 18
	pipeLoc.tlSecond.y = pipeLoc.brFirst.y + pipePieceUnit
	pipeLoc.brSecond.x = vMaxX - 15
	pipeLoc.brSecond.y = vMaxY - 1

	return pipeLoc
}

//
func drawPipe(v *gocui.View) {
	drawBlock(pipeLocations[0].tlFirst, pipeLocations[0].brFirst, v)
	drawBlock(pipeLocations[0].tlSecond, pipeLocations[0].brSecond, v)
}

// draw top and bottom borders
func drawBorders(v *gocui.View) error {

	// draw top and bottom edge borders
	drawBlock(point{x: 0, y: 0}, point{x: vMaxX - 1, y: 0}, v)
	drawBlock(point{x: 0, y: vMaxY - 1}, point{x: vMaxX - 1, y: vMaxY - 1}, v)

	return nil
}

// the function draws a block both horizontal and vertical
func drawBlock(p1 point, p2 point, v *gocui.View) {

	for y := p1.y; y <= p2.y; y++ {
		for x := p1.x; x <= p2.x; x++ {
			v.SetCursor(x, y)
			v.EditWrite('#')
		}
	}
}

func checkBottomCollision() bool {
	return birdPosition.y == vMaxY-1
}

func checkTopCollision() bool {
	return birdPosition.y == 0
}

func checkPipeCollisions() bool {
	return false
}

// silently quit
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
