package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/faiface/pixel/text"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

type Map struct {
	villagers []*Villager
	buildings []*Building

	houseCount int
}

var (
	focused   = false
	panelRect = pixel.R(defaultWidth/2*.1, (defaultHeight-400)/2, defaultWidth/2*.9, (defaultHeight-400)/2+400)

	houseButton   = pixel.R(.05, .55, .45, .95)
	labButton     = pixel.R(.05, .05, .45, .45)
	cantinaButton = pixel.R(.55, .55, .95, .95)
	repairButton  = pixel.R(.55, .05, .95, .45)

	landing = -1 // kind of building you want to land
)

type kindOfBuildings int

const (
	villagerSpeed = 50
	foodSupply    = .001

	house   = 0
	cantina = 1
	lab     = 2
	repair  = 3

	houseHalfSize = 30
	houseCost     = .05

	cantinaRad  = 40
	cantinaCost = .2

	labHalfSize = 50
	labCost     = .5
)

type Building struct {
	kind       kindOfBuildings
	position   pixel.Rect
	life       int
	creating   bool
	buildSince int

	data int
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
		life:     5,
		data:     0,
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
	case lab:
		if b.creating {
			imd.Color = colornames.Chocolate
		} else {
			imd.Color = colornames.Darkorchid
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

func (m *Map) clickBuilding() *Building {
	for _, b := range m.buildings {
		if b.position.Contains(mouseStart) {
			return b
		}
	}
	return nil
}

func (m *Map) canLandBuilding(target pixel.Rect) bool {
	for _, b := range m.buildings {
		if b.position.Intersect(target).Area() != 0 {
			return false
		}
	}
	return true
}

func (m *Map) forSelected(fn func(int, *Villager)) {
	for i, v := range m.villagers {
		if v.selected {
			fn(i, v)
		}
	}
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
	case lab:
		if pressed == 1 && !focused && p.scrap-labCost >= 0 {
			target := pixel.R(mouseLocation.X-labHalfSize, mouseLocation.Y-labHalfSize, mouseLocation.X+labHalfSize, mouseLocation.Y+labHalfSize)
			if target.Intersect(rightSide).Area() == target.Area() && m.canLandBuilding(target) {
				b := createBuilding(lab, target)
				m.buildings = append(m.buildings, b)
				p.scrap -= labCost
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
	case repair:
		if pressed == 1 && !focused {
			b := m.clickBuilding()
			if b != nil {
				m.forSelected(func(i int, v *Villager) {
					v.target = b
				})
			}
		}
	}

	for _, v := range m.villagers {
		p.food -= foodSupply

		// TODO: add inerty
		v.rigidBody.velocity = v.rigidBody.velocity.
			Add(p.rigidBody.velocity).
			Scaled(.1)

		if v.target != nil && v.target.life == 0 {
			v.target = nil
		}

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
					if v.target.kind == house {
						m.houseCount++
					}
					v.target.buildSince = t
					v.target = nil
				}
			}
		}

		v.rigidBody.physics(dt)
	}

	var aliveBuildings []*Building
	houses := 0
	for _, b := range m.buildings {
		if !b.creating {
			switch b.kind {
			case cantina:
				p.food += 0.01
				if p.food >= 1 {
					p.food = 1
				}
			case house:
				if (t+b.buildSince)%250 == 0 && len(m.villagers) < m.houseCount*5 {
					if b.data < 5 {
						m.villagers = append(m.villagers, &Villager{
							rigidBody: NewRigidBodyBySize(b.position.Center().X+rand.Float64()*houseHalfSize,
								b.position.Center().Y+rand.Float64()*houseHalfSize, 10, 10, pixel.ZV),
						})
						b.data++
					}
				}
			case lab:
				p.energy += 0.005
				if p.energy >= 1 {
					p.energy = 1
				}
			}
		}

		if p.isHit {
			b.life -= rand.Intn(2)
			if b.life < 0 {
				b.life = 0
			}
		}

		if b.life > 0 {
			if b.kind == house && !b.creating {
				houses++
			}
			aliveBuildings = append(aliveBuildings, b)
		}
	}
	m.buildings = aliveBuildings
	m.houseCount = houses

	if p.food <= 0 {
		p.food = 0

		if t%200 == 0 {
			var survivingVillager []*Villager
			for i, v := range m.villagers {
				if i+1 == rand.Intn(len(m.villagers)) {
					survivingVillager = append(survivingVillager, v)
				}
			}

			m.villagers = survivingVillager
		}
	}

	if len(m.villagers) == 0 || len(m.villagers) >= 200 {
		// end game
		screen = endScreen
	}

	panelRect = pixel.R(width/2+width/2*.1, (height-400)/2, width/2+width/2*.9, (height-400)/2+400)
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

			if rect := adapt(labButton, panelRect); rect.Contains(mousePosition) {
				landing = lab
				focused = false
			}

			if rect := adapt(repairButton, panelRect); rect.Contains(mousePosition) {
				landing = repair
				focused = false
			}

		}
	}
}

// Graphic functions

func adapt(rect1, rect2 pixel.Rect) pixel.Rect {
	return pixel.R(rect1.Min.X*rect2.W(),
		rect1.Min.Y*rect2.H(),
		rect1.Max.X*rect2.W(),
		rect1.Max.Y*rect2.H()).Moved(rect2.Min)
}

func drawButton(imd *imdraw.IMDraw, txt string, textScale float64, btn, ref pixel.Rect) {
	rect := adapt(btn, ref)
	imd.Color = colornames.Black
	imd.Push(rect.Min)
	imd.Push(rect.Max)
	imd.Rectangle(1)

	label := text.New(rect.Center(), uiFont)
	label.Color = color.Black
	label.Dot.X -= label.BoundsOf(txt).W() / 2
	fmt.Fprintf(label, txt)
	label.Draw(canvas, pixel.IM.Scaled(label.Orig, textScale))
}

func drawPanel(imd *imdraw.IMDraw) {
	// container
	imd.Color = colornames.Wheat
	imd.Push(panelRect.Min)
	imd.Push(panelRect.Max)
	imd.Rectangle(0)

	// house button
	drawButton(imd, "Build house", 1.1, houseButton, panelRect)
	drawButton(imd, "Build cantina", 1.1, cantinaButton, panelRect)
	drawButton(imd, "Build lab", 1.1, labButton, panelRect)
	drawButton(imd, "Repair building", 1.1, repairButton, panelRect)
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
	case lab:
		imag.Color = color.RGBA{255, 0, 0, 100}
		imag.Push(mouseLocation.Add(pixel.V(-labHalfSize, -labHalfSize)))
		imag.Push(mouseLocation.Add(pixel.V(labHalfSize, labHalfSize)))
		imag.Rectangle(0)
	case cantina:
		imag.Color = color.RGBA{255, 0, 0, 100}
		imag.Push(mouseLocation)
		imag.Circle(cantinaRad, 0)
	}

	txt := fmt.Sprintf("Population: %d/%d", len(m.villagers), m.houseCount*5)
	label := text.New(pixel.V(width/2+20, height-4*gap-4*h), uiFont)
	label.Color = color.Black
	fmt.Fprintf(label, txt)
	label.Draw(canvas, pixel.IM)

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
