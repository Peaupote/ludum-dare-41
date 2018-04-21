package main

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

type Map struct {
	villagers []*Villager
	buildings []*Buildings
}

var focused = false

var panelRect = pixel.R(defaultWidth/2+100, 100, defaultWidth-100, defaultHeight-100)

type kindOfBuildings int

type Buildings struct {
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

func mouseRect(v1, v2 pixel.Vec) pixel.Rect {
	if v1.X < v2.X && v1.Y < v2.Y {
		return pixel.R(v1.X, v1.Y, v2.X, v2.Y)
	} else if v1.X > v2.X && v1.Y > v2.Y {
		return pixel.R(v2.X, v2.Y, v1.X, v1.Y)
	} else if v1.X > v2.X {
		return pixel.R(v2.X, v1.Y, v1.X, v2.Y)
	}
	return pixel.R(v1.X, v2.Y, v2.X, v1.Y)
}

func getSelected(villagers []*Villager) []*Villager {
	var selected []*Villager
	rect := mouseRect(mouseStart, mousePosition).Intersect(rightSide)
	for _, v := range villagers {
		if v.rigidBody.hit(rect) {
			selected = append(selected, v)
		}
	}
	return selected
}

func drawPanel(imd *imdraw.IMDraw) {
	imd.Color = colornames.Wheat
	imd.Push(panelRect.Min)
	imd.Push(panelRect.Max)
	imd.Rectangle(0)
}

func (m *Map) update(dt float64) {
	panelRect = pixel.R(width/2+100, 100, width-100, height-100)
	if selected := getSelected(m.villagers); !pressed && len(selected) > 0 {
		focused = true
	}

	if focused && pressed && !panelRect.Contains(mouseStart) {
		focused = false
	}
}

func (m *Map) draw(imag *imdraw.IMDraw) {
	for _, v := range m.villagers {
		v.draw(imag)
	}

	if focused {
		drawPanel(imag)
	}

	m.drawSelectionZone(imag)
}

func (m *Map) drawSelectionZone(imag *imdraw.IMDraw) {
	if pressed {
		imag.Color = color.RGBA{0, 0, 1, 100}
		rect := mouseRect(mouseStart, mousePosition).Intersect(rightSide)
		imag.Push(rect.Min)
		imag.Push(rect.Max)
		imag.Rectangle(0)
	}
}
