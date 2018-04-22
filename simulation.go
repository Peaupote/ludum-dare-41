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

var (
	focused   = false
	panelRect = pixel.R(defaultWidth/2+100, 100, defaultWidth-100, defaultHeight-100)

	houseButton   = pixel.R(.05, .55, .45, .95)
	cantinaButton = pixel.R(.55, .55, .95, .95)

	landing = -1 // kind of building you want to land
)

type kindOfBuildings int

const (
	villagerSpeed = 50
	foodSupply    = .001

	house   = 0
	cantina = 1

	houseHalfSize = 30
	houseCost     = .05

	cantinaRad  = 40
	cantinaCost = .50
)

type Building struct {
	kind     kindOfBuildings
	position pixel.Rect
	life     int
	creating bool
}

type Villager struct {
	rigidBody *RigidBody
	target    *Building
	selected  bool
}

func (v *Villager) draw(imag *imdraw.IMDraw) {
	imag.Color = colornames.Blue
	v.rigidBody.draw(imag)
}

func createBuilding(k kindOfBuildings, p pixel.Rect) *Building {
	return &Building{
		kind:     k,
		position: p,
		creating: true,
		life:     0,
	}
}

func (m *Map) setTargetForSelectedVillagers(b *Building) {
	for _, v := range m.villagers {
		if v.selected {
			v.target = b
		}
	}
}

func (b *Building) draw(imd *imdraw.IMDraw) {
	switch b.kind {
	case house:
		if b.creating {
			imd.Color = colornames.Darkkhaki
		} else {
			imd.Color = colornames.Brown
		}

		imd.Push(b.position.Min)
		imd.Push(b.position.Max)
		imd.Rectangle(0)
	case cantina:
		if b.creating {
			imd.Color = colornames.Greenyellow
		} else {
			imd.Color = colornames.Green
		}
		imd.Push(b.position.Center())
		imd.Circle(cantinaRad, 0)
	}

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

func (m *Map) canLandBuilding(target pixel.Rect) bool {
	for _, b := range m.buildings {
		if b.position.Intersect(target).Area() != 0 {
			return false
		}
	}
	return true
}

func (m *Map) update(dt float64, p *Player) {
	if rightPressed > 0 || escape > 0 {
		landing = -1
		for _, v := range m.villagers {
			v.selected = false
		}
	}

	switch landing {
	case house:
		if pressed == 1 && !focused && p.scrap-houseCost >= 0 {
			target := pixel.R(mouseLocation.X-houseHalfSize, mouseLocation.Y-houseHalfSize, mouseLocation.X+houseHalfSize, mouseLocation.Y+houseHalfSize)
			if target.Intersect(rightSide).Area() == target.Area() && m.canLandBuilding(target) {
				b := createBuilding(house, target)
				m.buildings = append(m.buildings, b)
				p.scrap -= houseCost
				m.setTargetForSelectedVillagers(b)
			}
		}
	case cantina:
		if pressed == 1 && !focused && p.scrap-cantinaCost >= 0 {
			target := pixel.R(mouseLocation.X-cantinaRad, mouseLocation.Y-cantinaRad, mouseLocation.X+cantinaRad, mouseLocation.Y+cantinaRad)
			if target.Intersect(rightSide).Area() == target.Area() && m.canLandBuilding(target) {
				b := createBuilding(cantina, target)
				m.buildings = append(m.buildings, b)
				p.scrap -= cantinaCost
				m.setTargetForSelectedVillagers(b)
			}
		}
	}

	for _, v := range m.villagers {
		p.food -= foodSupply

		// TODO: add inerty
		v.rigidBody.velocity = v.rigidBody.velocity.
			Add(p.rigidBody.velocity).
			Scaled(.1)

		if v.target != nil {
			v.rigidBody.velocity = v.rigidBody.velocity.
				Add(v.target.position.Center().
					Add(v.rigidBody.body.Center().Scaled(-1)).
					Scaled(villagerSpeed * dt))

			if v.rigidBody.body.Intersect(v.target.position).Area() == v.rigidBody.body.Area() {
				v.rigidBody.velocity = pixel.ZV
				v.target.life++

				if v.target.life >= 100 {
					v.target.life = 100
					v.target.creating = false
					v.target = nil
				}
			}
		}

		v.rigidBody.physics(dt)
	}

	for _, b := range m.buildings {
		if !b.creating {
			switch b.kind {
			case cantina:
				p.food += 0.01
				if p.food >= 1 {
					p.food = 1
				}
			case house:
				if len(m.villagers) < 10 {
					m.villagers = append(m.villagers, &Villager{
						rigidBody: NewRigidBodyBySize(b.position.Center().X, b.position.Center().Y, 10, 10, pixel.ZV),
					})
				}
			}
		}
	}

	if p.food <= 0 {
		// TODO: kill random villagers or lose
		p.food = 0
	}

	panelRect = pixel.R(width/2+100, 100, width-100, height-100)
	if !focused {
		if selected := getSelected(m.villagers); pressed == 0 && len(selected) > 0 {
			focused = true
			for _, v := range selected {
				v.selected = true
			}
		}
	}

	if focused && pressed == 1 {
		if !panelRect.Contains(mouseStart) {
			focused = false
		} else {
			if rect := adapt(houseButton, panelRect); rect.Contains(mousePosition) {
				landing = house
				focused = false
			}

			if rect := adapt(cantinaButton, panelRect); rect.Contains(mousePosition) {
				landing = cantina
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
		imag.Push(mouseLocation.Add(pixel.V(-houseHalfSize, -houseHalfSize)))
		imag.Push(mouseLocation.Add(pixel.V(houseHalfSize, houseHalfSize)))
		imag.Rectangle(0)
	case cantina:
		imag.Color = color.RGBA{255, 0, 0, 100}
		imag.Push(mouseLocation)
		imag.Circle(cantinaRad, 0)
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