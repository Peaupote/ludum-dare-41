package main

import (
	"fmt"

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
)

func menuUpdate(dt float64, win *pixelgl.Window, imd *imdraw.IMDraw) {
	if pressed > 0 {
		if rect := adapt(menuButton, canvas.Bounds()); rect.Contains(mouseStart) {
			startGame()
			screen = gameScreen
		}

		if rect := adapt(quitButton, canvas.Bounds()); rect.Contains(mouseStart) {
			quit = true
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
}
