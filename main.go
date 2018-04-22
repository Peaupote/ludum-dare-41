package main

import (
	"image"
	"os"
	"time"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
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
	escape = 0

	rightPressed  = 0
	pressed       = 0
	mouseStart    pixel.Vec
	mousePosition pixel.Vec
	mouseLocation pixel.Vec

	leftSide  pixel.Rect
	rightSide pixel.Rect

	uiFont = text.NewAtlas(basicfont.Face7x13, text.ASCII)

	width  float64
	height float64
	t      int
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

	if win.JustReleased(pixelgl.KeyEscape) {
		escape = 0
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

	if win.Pressed(pixelgl.KeyEscape) {
		escape++
	}

	//

	mouseLocation = win.MousePosition()

	if win.JustPressed(pixelgl.MouseButtonLeft) {
		mouseStart = win.MousePosition()
	}

	if win.Pressed(pixelgl.MouseButtonLeft) {
		pressed++
	}

	if win.Pressed(pixelgl.MouseButtonLeft) {
		mousePosition = win.MousePosition()
	}

	if win.JustReleased(pixelgl.MouseButtonLeft) {
		pressed = 0
	}

	if win.Pressed(pixelgl.MouseButtonRight) {
		rightPressed++
	}

	if win.JustReleased(pixelgl.MouseButtonRight) {
		rightPressed = 0
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

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return pixel.PictureDataFromImage(img), nil
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
	t = 0

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

	spaceBackground, err := loadPicture("./assets/space-bkg.png")
	if err != nil {
		panic(err)
	}

	sprite := pixel.NewSprite(spaceBackground, spaceBackground.Bounds())

	var ovnis []*Ovni

	vils := []*Villager{&Villager{
		rigidBody: NewRigidBodyBySize(defaultWidth/2+100, 100, 10, 10, pixel.ZV),
	}}

	m := &Map{
		villagers: vils,
		buildings: nil,
	}

	win.SetSmooth(true)

	for !win.Closed() {
		if win.JustPressed(pixelgl.KeyQ) {
			return
		}

		dt := time.Since(last).Seconds()
		t++
		last = time.Now()

		applyControls(win)

		if enter > 0 && player.energy > .01 {
			dt /= 3
			player.energy -= .01
			if player.energy < 0 {
				player.energy = 0
			}
		}

		// update
		ovnis = updateUniverse(dt, ovnis)
		ovnis = player.upadte(dt, ovnis)
		m.update(dt, player)

		win.Clear(colornames.Skyblue)
		imd.Clear()

		width = win.Bounds().W()
		height = win.Bounds().H()

		// TODO: clean up here
		leftSide = pixel.R(0, 0, width/2, height)
		rightSide = pixel.R(width/2, 0, width, height)
		sprite.Draw(win, pixel.IM.
			Moved(leftSide.Center()).
			ScaledXY(leftSide.Center(), pixel.V(width/(2*spaceBackground.Bounds().W()), 1)))
		middleBar(imd)

		for _, o := range ovnis {
			o.draw(imd)
		}

		m.draw(imd)
		player.draw(imd)
		imd.Draw(win)

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
