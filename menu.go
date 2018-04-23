package main

import (
	"fmt"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

var (
	title     = "{insert title here}"
	menuLines = []string{
		"Hi ! You are the commander",
		"of a spaceship travelling",
		"through intergallactic",
		"spaces. You must control",
		"the ship but also manage",
		"workers' lives in the ship.",
		"Build houses to raise more",
		"workers, cantina to provide",
		"them food and lab to increase",
		"your energy production",
		"",
		"Good luck !",
	}

	loseLines = []string{
		"To bad...",
		"All your crew is dead",
		"It's the end of the adventure",
	}

	winLines = []string{
		"Well done !",
		"You're the best commander",
		"i've ever met",
	}

	controlLines = []string{
		"Arrows - Move",
		"Space  - Shoot",
		"Tab    - Switch weapon",
		"Enter  - Bullet time",
		"",
		"You can select workers",
		"by drawing a square",
		"including them with your mouse",
	}

	menuButton = pixel.R(0.25, 0.1, .45, .2)
	quitButton = pixel.R(0.55, 0.1, .75, .2)

	hardMode bool
)

func menuUpdate(dt float64, win *pixelgl.Window, imd *imdraw.IMDraw) {
	if pressed == 1 {
		if rect := adapt(menuButton, canvas.Bounds()); rect.Contains(mouseStart) {
			startGame()
			screen = gameScreen
		}

		if rect := adapt(quitButton, canvas.Bounds()); rect.Contains(mouseStart) {
			quit = true
		}

		if btn := adapt(menuButton, canvas.Bounds()); pixel.R(btn.Min.X, btn.Min.Y-20, btn.Min.X+100, btn.Min.Y-10).Contains(mouseStart) {
			hardMode = !hardMode
		}
	}
}

func endUpdate(dt float64, win *pixelgl.Window, imd *imdraw.IMDraw) {
	menuUpdate(dt, win, imd)
}

func menuRender(dt float64, win *pixelgl.Window, imd *imdraw.IMDraw) {
	drawMenu(title, menuLines, "Start play", win, imd)
}

func endRender(dt float64, win *pixelgl.Window, imd *imdraw.IMDraw) {
	endHeader := "Loser !"
	if len(m.villagers) >= popToWin {
		endHeader = "Whoa you made it !"
	}

	lines := loseLines
	if len(m.villagers) >= popToWin {
		lines = winLines
	}

	drawMenu(endHeader, lines, "Play again", win, imd)
}

func drawMenu(headerTxt string, txt []string, buttonText string, win *pixelgl.Window, imd *imdraw.IMDraw) {
	header := text.New(canvas.Bounds().Center().Add(pixel.V(0, width/4)), uiFont)
	header.Dot.X -= header.BoundsOf(headerTxt).W() / 2
	fmt.Fprintf(header, headerTxt)
	header.Draw(canvas, pixel.IM.Scaled(header.Orig, 5))

	label := text.New(pixel.V(width*.05, height*.65), uiFont)
	fmt.Fprintf(label, "Controls")
	label.Draw(canvas, pixel.IM.Scaled(label.Orig, 3))

	label = text.New(pixel.V(width*.05, height*.6), uiFont)
	for _, line := range controlLines {
		fmt.Fprintln(label, line)
	}
	label.Draw(canvas, pixel.IM.Scaled(label.Orig, 2))

	label = text.New(pixel.V(width*.55, height*.65), uiFont)
	for _, line := range txt {
		fmt.Fprintln(label, line)
	}
	label.Draw(canvas, pixel.IM.Scaled(label.Orig, 2))

	drawButton(imd, buttonText, 2, menuButton, canvas.Bounds())
	drawButton(imd, "Quit", 2, quitButton, canvas.Bounds())

	btn := adapt(menuButton, canvas.Bounds())
	if hardMode {
		imd.Color = colornames.Green
		imd.Push(btn.Min.Add(pixel.V(0, -20)))
		imd.Push(btn.Min.Add(pixel.V(10, -10)))
		imd.Rectangle(0)
	}

	imd.Color = colornames.Grey
	imd.Push(btn.Min.Add(pixel.V(0, -20)))
	imd.Push(btn.Min.Add(pixel.V(10, -10)))
	imd.Rectangle(1)

	hard := text.New(btn.Min.Add(pixel.V(15, -20)), uiFont)
	fmt.Fprintf(hard, "Hard mode (must reac 150 population to win)")
	hard.Color = colornames.Black
	hard.Draw(canvas, pixel.IM)
}
