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
		"Move the dxcfvgbhjklmù",
		"rfghjklkhgfdxgkolpmgtfdrtgyuio",
		"cfvghjklmlkjhgfdghjkl",
		"ghjkl",
	}
	loseLines = []string{
		"You are a loser",
		"fdkhfdjkhf %s might be cool to insert text here",
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

	drawMenu(endHeader, loseLines, "Play again", win, imd)
}

func drawMenu(headerTxt string, txt []string, buttonText string, win *pixelgl.Window, imd *imdraw.IMDraw) {
	header := text.New(canvas.Bounds().Center().Add(pixel.V(0, width/4)), uiFont)
	header.Dot.X -= header.BoundsOf(headerTxt).W() / 2
	fmt.Fprintf(header, headerTxt)
	header.Draw(canvas, pixel.IM.Scaled(header.Orig, 5))

	label := text.New(canvas.Bounds().Center(), uiFont)
	for _, line := range txt {
		label.Dot.X -= label.BoundsOf(line).W() / 2
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
