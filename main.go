package main

import (
	"encoding/csv"
	"image"
	"io"
	"os"
	"strconv"
	"time"

	"image/color"
	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/pkg/errors"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

const (
	defaultWidth  = 1024
	defaultHeight = 750

	menuScreen = 0
	gameScreen = 1
	endScreen  = 2
)

var (
	screen = menuScreen

	left   = 0
	right  = 0
	bottom = 0
	top    = 0
	space  = 0
	enter  = 0
	tab    = 0
	escape = 0

	quit = false

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

	canvas    *pixelgl.Canvas
	topCanvas *pixelgl.Canvas

	player          *Player
	m               *Map
	ovnis           []*Ovni
	sprite          *pixel.Sprite
	spaceBackground pixel.Picture
)

func loadAnimationSheet(sheetPath, descPath string, frameWidth float64) (sheet pixel.Picture, anims map[string][]pixel.Rect, err error) {
	// total hack, nicely format the error at the end, so I don't have to type it every time
	defer func() {
		if err != nil {
			err = errors.Wrap(err, "error loading animation sheet")
		}
	}()

	// open and load the spritesheet
	sheetFile, err := os.Open(sheetPath)
	if err != nil {
		return nil, nil, err
	}
	defer sheetFile.Close()
	sheetImg, _, err := image.Decode(sheetFile)
	if err != nil {
		return nil, nil, err
	}
	sheet = pixel.PictureDataFromImage(sheetImg)

	// create a slice of frames inside the spritesheet
	var frames []pixel.Rect
	for x := 0.0; x+frameWidth <= sheet.Bounds().Max.X; x += frameWidth {
		frames = append(frames, pixel.R(
			x,
			0,
			x+frameWidth,
			sheet.Bounds().H(),
		))
	}

	descFile, err := os.Open(descPath)
	if err != nil {
		return nil, nil, err
	}
	defer descFile.Close()

	anims = make(map[string][]pixel.Rect)

	// load the animation information, name and interval inside the spritesheet
	desc := csv.NewReader(descFile)
	for {
		anim, err := desc.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}

		name := anim[0]
		start, _ := strconv.Atoi(anim[1])
		end, _ := strconv.Atoi(anim[2])

		anims[name] = frames[start : end+1]
	}

	return sheet, anims, nil
}

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

	canvas = pixelgl.NewCanvas(win.Bounds())
	topCanvas = pixelgl.NewCanvas(win.Bounds())
	frameCanvas := pixelgl.NewCanvas(win.Bounds())

	win.SetSmooth(true)

	for !win.Closed() {
		if quit || win.JustPressed(pixelgl.KeyQ) {
			return
		}

		dt := time.Since(last).Seconds()
		t++
		last = time.Now()

		applyControls(win)

		switch screen {
		case menuScreen:
			menuUpdate(dt, win, imd)
		case gameScreen:
			gameUpdate(dt, win, imd)
		case endScreen:
			endUpdate(dt, win, imd)
		}

		topCanvas.Clear(color.RGBA{0, 0, 0, 0})
		frameCanvas.Clear(color.RGBA{0, 0, 0, 0})
		canvas.Clear(color.RGBA{0, 0, 0, 0})
		win.Clear(colornames.Skyblue)
		imd.Clear()

		width = win.Bounds().W()
		height = win.Bounds().H()

		canvas.SetBounds(win.Bounds())
		switch screen {
		case menuScreen:
			menuRender(dt, win, imd)
		case gameScreen:
			gameRender(dt, win, imd)
		case endScreen:
			endRender(dt, win, imd)
		}

		canvas.Draw(frameCanvas, pixel.IM.Moved(win.Bounds().Center()))
		imd.Draw(frameCanvas)
		topCanvas.Draw(frameCanvas, pixel.IM.Moved(win.Bounds().Center()))
		frameCanvas.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		win.Update()
	}
}

func startGame() {
	sheet, anims, err := loadAnimationSheet("./assets/ship.png", "./assets/ship.csv", 30)
	if err != nil {
		panic(err)
	}

	if hardMode {
		popToWin = 150
	} else {
		popToWin = 50
	}

	player = &Player{
		rigidBody: &RigidBody{
			body:     pixel.R(200, 200, 260, 300),
			velocity: pixel.ZV,
		},
		mode:   shootLaser,
		energy: 1,
		food:   1,
		scrap:  1,

		sheet: sheet,
		anims: anims,
		rate:  1.0 / 10,
	}

	sp, err := loadPicture("./assets/space-bkg.png")
	if err != nil {
		panic(err)
	}

	spaceBackground = sp

	sprite = pixel.NewSprite(spaceBackground, spaceBackground.Bounds())

	icons, err := loadPicture("./assets/shootMode.png")
	if err != nil {
		panic(err)
	}
	shootModeSprite = pixel.NewSprite(icons, icons.Bounds())

	vils := []*Villager{NewVillager(defaultWidth/2+100, 100)}
	bils := []*Building{
		createBuilding(cantina, pixel.R(642, 69, 722, 149)),
		createBuilding(house, pixel.R(592, 177, 652, 237)),
	}

	bils[0].creating = false
	bils[0].life = 100.0

	bils[1].creating = false
	bils[1].life = 100.0

	ovnis = []*Ovni{}
	landing = -1
	focused = false

	m = &Map{
		villagers: vils,
		buildings: bils,
	}
}

func gameUpdate(dt float64, win *pixelgl.Window, imd *imdraw.IMDraw) {
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
}

func gameRender(dt float64, win *pixelgl.Window, imd *imdraw.IMDraw) {
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
}

func main() {
	pixelgl.Run(run)
}
