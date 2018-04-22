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
	houseCost     = .15

	cantinaRad  = 40
	cantinaCost = .2

	labHalfSize = 50
	labCost     = .5

	popToWin = 10
)

type Building struct {
	kind       kindOfBuildings
	position   pixel.Rect
	life       float64
	creating   bool
	buildSince int

	sheet   pixel.Picture
	anims   map[string][]pixel.Rect
	rate    float64
	counter float64
	index   int

	frame pixel.Rect

	sprite *pixel.Sprite
}

type Villager struct {
	rigidBody *RigidBody
	target    *Building
	selected  bool

	sheet   pixel.Picture
	anims   map[string][]pixel.Rect
	rate    float64
	counter float64
	index   int

	frame pixel.Rect

	sprite *pixel.Sprite
}

func NewVillager(x, y float64) *Villager {
	sheet, anims, err := loadAnimationSheet("./assets/villager.png", "./assets/villager.csv", 10)
	if err != nil {
		panic(err)
	}

	return &Villager{
		rigidBody: NewRigidBodyBySize(x, y, 10, 10, pixel.ZV),
		sheet:     sheet,
		anims:     anims,
		rate:      1.0 / 10,
	}
}

func (v *Villager) draw(imag *imdraw.IMDraw) {
	if v.sprite == nil {
		v.sprite = pixel.NewSprite(nil, pixel.Rect{})
	}

	v.sprite.Set(v.sheet, v.frame)
	v.sprite.Draw(canvas, pixel.IM.
		Moved(v.rigidBody.body.Center()).
		Scaled(v.rigidBody.body.Center(), 2))
}

func createBuilding(k kindOfBuildings, p pixel.Rect) *Building {
	var sheet pixel.Picture
	var anims map[string][]pixel.Rect
	switch k {
	case house:
		s, a, err := loadAnimationSheet("./assets/house.png", "./assets/building.csv", 60)
		if err != nil {
			panic(err)
		}
		sheet = s
		anims = a
	case lab:
		s, a, err := loadAnimationSheet("./assets/lab.png", "./assets/lab.csv", 100)
		if err != nil {
			panic(err)
		}
		sheet = s
		anims = a
	}

	return &Building{
		kind:     k,
		position: p,
		creating: true,
		life:     5,

		sheet: sheet,
		anims: anims,
		rate:  1.0 / 10,
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
	case house, lab:
		if b.sprite == nil {
			b.sprite = pixel.NewSprite(nil, pixel.Rect{})
		}

		b.sprite.Set(b.sheet, b.frame)
		b.sprite.Draw(canvas, pixel.IM.
			Moved(b.position.Center()))
	case cantina:
		if b.creating {
			imd.Color = colornames.Greenyellow
		} else {
			imd.Color = colornames.Green
		}
		imd.Push(b.position.Center())
		imd.Circle(cantinaRad, 0)
	}

	txt := fmt.Sprintf("Life: %d/100", int(b.life))
	label := text.New(pixel.V(b.position.Center().X, b.position.Min.Y-10), uiFont)
	label.Dot.X -= label.BoundsOf(txt).W() / 2
	label.Color = color.Black
	fmt.Fprintf(label, txt)
	label.Draw(canvas, pixel.IM)

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

func (m *Map) clearSelected() {
	for _, v := range m.villagers {
		v.selected = false
	}
}

func (m *Map) update(dt float64, p *Player) {
	if rightPressed > 0 || escape > 0 {
		focused = false
		landing = -1
		m.clearSelected()
		mousePosition = pixel.ZV
		mouseStart = pixel.ZV
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
			landing = -1
			m.clearSelected()
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
			m.clearSelected()
			landing = -1
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
			m.clearSelected()
			landing = -1
		}
	case repair:
		if pressed == 1 && !focused {
			b := m.clickBuilding()
			if b != nil {
				m.forSelected(func(i int, v *Villager) {
					v.target = b
				})
			}
			landing = -1
			m.clearSelected()
		}
	}

	for _, v := range m.villagers {
		p.food -= foodSupply

		v.counter++
		if v.counter > 20 {
			v.index = 1 - v.index
			v.counter = 0
		}
		v.frame = v.anims["Idle"][v.index]

		r := v.rigidBody

		half := width/2 + 20
		if r.body.Min.X < half {
			r.body = r.body.Moved(pixel.V(half-r.body.Min.X, 0))
			r.velocity = pixel.ZV
		}

		if r.body.Max.X > width {
			r.body = r.body.Moved(pixel.V(width-r.body.Max.X, 0))
			r.velocity = pixel.ZV
		}

		if r.body.Min.Y < 0 {
			r.body = r.body.Moved(pixel.V(0, -r.body.Min.Y))
			r.velocity = pixel.ZV
		}

		if r.body.Max.Y > height {
			r.body = r.body.Moved(pixel.V(0, height-r.body.Max.Y))
			r.velocity = pixel.ZV
		}

		// TODO: add inerty
		r.velocity = r.velocity.
			Add(p.rigidBody.velocity).Scaled(.5)

		if v.target != nil && v.target.life == 0 {
			v.target = nil
		}

		if v.target != nil {
			r.velocity = r.velocity.
				Add(v.target.position.Center().
					Add(r.body.Center().Scaled(-1)).
					Scaled(villagerSpeed * dt))

			if r.body.Intersect(v.target.position).Area() == r.body.Area() {
				r.velocity = pixel.ZV
				v.target.life += .5

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

		r.physics(dt)
	}

	var aliveBuildings []*Building
	houses := 0
	for _, b := range m.buildings {
		if !b.creating {
			switch b.kind {
			case cantina:
				p.food += foodSupply * 10
				if p.food >= 1 {
					p.food = 1
				}
			case house:
				if (t+b.buildSince)%300 == 0 && len(m.villagers) < m.houseCount*5 && p.food-0.05 > 0 {
					m.villagers = append(m.villagers, NewVillager(
						b.position.Center().X+rand.Float64()*houseHalfSize,
						b.position.Center().Y+rand.Float64()*houseHalfSize))
					p.food -= 0.05
					if p.food < 0 {
						p.food = 0
					}
				}
			case lab:
				p.energy += 0.0005
				if p.energy >= 1 {
					p.energy = 1
				}
			}
		}

		if p.isHit {
			b.life -= float64(rand.Intn(3))
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

		state := "Good"
		if b.life < 70 {
			state = "Med"
		}

		if b.life < 30 {
			state = "Bad"
		}

		if b.kind != cantina {
			b.counter++
			if b.counter > 50 {
				b.counter = 0
				b.index = (b.index + 1) % len(b.anims[state])
			}
			b.frame = b.anims[state][b.index]
		}

	}
	m.buildings = aliveBuildings
	m.houseCount = houses

	if p.food <= 0 {
		p.food = 0

		if t%500 == 0 {
			var survivingVillager []*Villager
			r := rand.Intn(len(m.villagers))
			for i, v := range m.villagers {
				if i+1 != r {
					survivingVillager = append(survivingVillager, v)
				}
			}

			m.villagers = survivingVillager
		}
	}

	if len(m.villagers) == 0 || len(m.villagers) >= popToWin {
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
			m.clearSelected()
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
	label.Draw(topCanvas, pixel.IM.Scaled(label.Orig, textScale))
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

	txt := fmt.Sprintf("press ESC to close this panel")
	label := text.New(panelRect.Center().Add(pixel.V(0, 5-panelRect.H()/2)), uiFont)
	label.Color = color.Black
	label.Dot.X -= label.BoundsOf(txt).W() / 2
	fmt.Fprintf(label, txt)
	label.Draw(topCanvas, pixel.IM)
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
	label.Draw(topCanvas, pixel.IM)

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
