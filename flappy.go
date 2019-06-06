package main

import (
	"log"
	"math/rand"
	"os"
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

	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeySpace, gocui.ModNone, birdUp); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}

// main function of game
func launchGame(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("flappy", 0, 0, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		birdPosition = point{x: 25, y: 10}
		v.FgColor = gocui.ColorYellow
		vMaxX, vMaxY = v.Size()
		rs := rand.NewSource(int64(time.Now().Second()))
		rng = rand.New(rs)
		pipePieceUnit = (vMaxY - 2) / 9

		ticker := time.NewTicker(100 * time.Millisecond)
		go func() {
			for range ticker.C {
				advanceGame()
				refreshBoard(g, v)
			}
		}()
	}
	return nil
}

// take the game one step ahead
// bird down, pipes move etc.
func advanceGame() error {
	birdDown(nil, nil)
	floatPipes()
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
	isCollided := false
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("flappy")
		if err != nil {

		}
		if checkCollision() {
			v.FgColor = gocui.ColorRed
			isCollided = true
		}
		v.Clear()
		drawBird(birdPosition, v)
		drawBorders(v)
		drawPipes(v)
		return nil
	})

	// wait for ui updating
	time.Sleep(50 * time.Millisecond)
	if isCollided {
		os.Exit(0)
	}
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

	pipeLoc.tlFirst.x = vMaxX - 5
	pipeLoc.tlFirst.y = 1
	pipeLoc.brFirst.x = vMaxX - 2
	pipeLoc.brFirst.y = pipeLoc.tlFirst.y + (firstPieceVolume * pipePieceUnit)

	pipeLoc.tlSecond.x = vMaxX - 5
	pipeLoc.tlSecond.y = pipeLoc.brFirst.y + pipePieceUnit*2
	pipeLoc.brSecond.x = vMaxX - 2
	pipeLoc.brSecond.y = vMaxY - 1

	return pipeLoc
}

//
func floatPipes() {
	for ix := range pipeLocations {
		pipeLocations[ix].tlFirst.x--
		pipeLocations[ix].brFirst.x--
		pipeLocations[ix].tlSecond.x--
		pipeLocations[ix].brSecond.x--
	}
}

//
func drawPipes(v *gocui.View) {

	if pipeLocations[1].tlFirst.x < vMaxX/2 {
		pipeLocations[0] = pipeLocations[1]
		pipeLocations[1] = calculateNewPipePosition()
	}

	for _, pLoc := range pipeLocations {
		drawPipe(pLoc, v)
	}

}

// draw pipe using pipe location
func drawPipe(pLoc pipeLocation, v *gocui.View) error {
	drawBlock(pLoc.tlFirst, pLoc.brFirst, v)
	drawBlock(pLoc.tlSecond, pLoc.brSecond, v)

	return nil
}

// draw top and bottom borders
func drawBorders(v *gocui.View) error {
	drawBlock(point{x: 0, y: 0}, point{x: vMaxX - 1, y: 0}, v)
	drawBlock(point{x: 0, y: vMaxY - 1}, point{x: vMaxX - 1, y: vMaxY - 1}, v)

	return nil
}

// the function draws a block both horizontal and vertical
func drawBlock(tl point, br point, v *gocui.View) {
	for y := tl.y; y <= br.y; y++ {
		for x := tl.x; x <= br.x; x++ {
			v.SetCursor(x, y)
			v.EditWrite('#')
		}
	}
}

func checkCollision() bool {
	return checkBottomCollision() ||
		checkTopCollision() ||
		checkPipeCollisions()
}

func checkBottomCollision() bool {
	return birdPosition.y > vMaxY-1
}

func checkTopCollision() bool {
	return birdPosition.y < 1
}

func checkPipeCollisions() bool {

	for _, pipeLoc := range pipeLocations {
		if (birdPosition.y < pipeLoc.brFirst.y ||
			birdPosition.y > pipeLoc.tlSecond.y) &&
			birdPosition.x == pipeLoc.brFirst.x {
			return true
		}
	}

	return false
}

// silently quit
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
