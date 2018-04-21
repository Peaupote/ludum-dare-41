package main

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

type Map struct {
	villagers []*Villager
	buildings []*Building
}

var focused = false

var (
	panelRect = pixel.R(defaultWidth/2+100, 100, defaultWidth-100, defaultHeight-100)

	houseButton   = pixel.R(.05, .55, .45, .95)
	cantinaButton = pixel.R(.55, .55, .95, .95)

	landing = -1 // kind of building you want to land
)

type kindOfBuildings int

const (
	house   = 0
	cantina = 1
)

type Building struct {
	kind     kindOfBuildings
	position pixel.Rect
	life     int
}

type Villager struct {
	rigidBody *RigidBody
}

func (v *Villager) draw(imag *imdraw.IMDraw) {
	imag.Color = colornames.Blue
	v.rigidBody.draw(imag)
}

func (b *Building) draw(imd *imdraw.IMDraw) {
	switch b.kind {
	case house:
		imd.Color = colornames.Darkkhaki
	}

	imd.Push(b.position.Min)
	imd.Push(b.position.Max)
	imd.Rectangle(0)
}

func getSelected(villagers []*Villager) []*Villager {
	var selected []*Villager
	rect := pixel.R(mouseStart.X, mouseStart.Y, mousePosition.X, mousePosition.Y).Norm().Intersect(rightSide)
	for _, v := range villagers {
		if v.rigidBody.hit(rect) {
			selected = append(selected, v)
		}
	}
	return selected
}

func adapt(rect1, rect2 pixel.Rect) pixel.Rect {
	return pixel.R(rect1.Min.X*rect2.W(),
		rect1.Min.Y*rect2.H(),
		rect1.Max.X*rect2.W(),
		rect1.Max.Y*rect2.H()).Moved(rect2.Min)
}

func drawPanel(imd *imdraw.IMDraw) {
	// container
	imd.Color = colornames.Wheat
	imd.Push(panelRect.Min)
	imd.Push(panelRect.Max)
	imd.Rectangle(0)

	// house button
	rect := adapt(houseButton, panelRect)
	imd.Color = colornames.Black
	imd.Push(rect.Min)
	imd.Push(rect.Max)
	imd.Rectangle(1)

	// label := text.New(pixel.V(0.1, .6).ScaledXY(pixel.V(panelRect.W(), panelRect.H()).Add(panelRect.Min)), uiFont)
	// fmt.Fprint(label, "Buy a house")
	// label.Draw(imd, pixel.IM)

	// cantina button
	rect = adapt(cantinaButton, panelRect)
	imd.Color = colornames.Black
	imd.Push(rect.Min)
	imd.Push(rect.Max)
	imd.Rectangle(1)

}

func (m *Map) update(dt float64) {
	if rightPressed > 0 {
		landing = 0
	}

	switch landing {
	case house:
		if pressed == 1 && !focused {
			m.buildings = append(m.buildings, &Building{
				kind:     house,
				position: pixel.R(mouseLocation.X-50, mouseLocation.Y-50, mouseLocation.X+50, mouseLocation.Y+50),
				life:     100,
			})
		}
	}

	panelRect = pixel.R(width/2+100, 100, width-100, height-100)
	if selected := getSelected(m.villagers); pressed == 0 && len(selected) > 0 {
		focused = true
	}

	if focused && pressed == 1 {
		if !panelRect.Contains(mouseStart) {
			focused = false
		} else {
			rect := adapt(houseButton, panelRect)
			if rect.Contains(mousePosition) {
				landing = house
				focused = false
			}
		}
	}
}

func (m *Map) draw(imag *imdraw.IMDraw) {
	for _, b := range m.buildings {
		b.draw(imag)
	}

	for _, v := range m.villagers {
		v.draw(imag)
	}

	if focused {
		drawPanel(imag)
	}

	switch landing {
	case house:
		imag.Color = color.RGBA{255, 0, 0, 100}
		// TODO: check for collisions
		imag.Push(mouseLocation.Add(pixel.V(-50, -50)))
		imag.Push(mouseLocation.Add(pixel.V(50, 50)))
		imag.Rectangle(0)
	}

	m.drawSelectionZone(imag)
}

func (m *Map) drawSelectionZone(imag *imdraw.IMDraw) {
	if pressed > 0 {
		imag.Color = color.RGBA{0, 0, 1, 100}
		rect := pixel.R(mouseStart.X, mouseStart.Y, mousePosition.X, mousePosition.Y).Norm().Intersect(rightSide)
		imag.Push(rect.Min)
		imag.Push(rect.Max)
		imag.Rectangle(0)
	}
}
