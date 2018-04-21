package main

import (
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	defaultWidth  = 1024
	defaultHeight = 750
)

var (
	left   = 0
	right  = 0
	bottom = 0
	top    = 0
	space  = 0
	enter  = 0
	tab    = 0

	width  float64
	height float64
)

func applyControls(win *pixelgl.Window) {
	if win.JustReleased(pixelgl.KeyLeft) {
		left = 0
	}

	if win.JustReleased(pixelgl.KeyRight) {
		right = 0
	}

	if win.JustReleased(pixelgl.KeyDown) {
		bottom = 0
	}

	if win.JustReleased(pixelgl.KeyUp) {
		top = 0
	}

	if win.JustReleased(pixelgl.KeySpace) {
		space = 0
	}

	if win.JustReleased(pixelgl.KeyEnter) {
		enter = 0
	}

	if win.JustReleased(pixelgl.KeyTab) {
		tab = 0
	}

	//

	if win.Pressed(pixelgl.KeyLeft) {
		left++
	}

	if win.Pressed(pixelgl.KeyRight) {
		right++
	}

	if win.Pressed(pixelgl.KeyDown) {
		bottom++
	}

	if win.Pressed(pixelgl.KeyUp) {
		top++
	}

	if win.Pressed(pixelgl.KeySpace) {
		space++
	}

	if win.Pressed(pixelgl.KeyEnter) {
		enter++
	}

	if win.Pressed(pixelgl.KeyTab) {
		tab++
	}
}

func middleBar(imd *imdraw.IMDraw) {
	// dual screen separation
	x := width / 2

	imd.Color = colornames.Black
	imd.Push(pixel.V(x, height))
	imd.Push(pixel.V(x+10, 0))
	imd.Rectangle(0)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Ludum dare 41 - {enter the name here}",
		Bounds: pixel.R(0, 0, defaultWidth, defaultHeight),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)
	last := time.Now()

	player := &Player{
		rigidBody: &RigidBody{
			body:     pixel.R(200, 300, 300, 200),
			velocity: pixel.ZV,
		},
		mode:   shootLaser,
		energy: .5,
		food:   0.5,
		scrap:  0.5,
	}

	var ovnis []*Ovni

	for !win.Closed() {
		if win.JustPressed(pixelgl.KeyEscape) {
			return
		}

		dt := time.Since(last).Seconds()
		last = time.Now()

		applyControls(win)

		// update
		ovnis = updateUniverse(dt, ovnis)
		ovnis = player.upadte(dt, ovnis)

		win.Clear(colornames.Skyblue)
		imd.Clear()

		width = win.Bounds().W()
		height = win.Bounds().H()
		middleBar(imd)

		for _, o := range ovnis {
			o.draw(imd)
		}
		player.draw(imd)
		imd.Draw(win)

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
