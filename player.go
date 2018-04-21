package main

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

type shootMode int

const (
	shootLaser   = 0
	shootBullets = 1
)

// Bullet is a single bullet the player can shoot
type Bullet struct {
	rigidBody *RigidBody
}

// Player represents the player
type Player struct {
	rigidBody *RigidBody
	mode      shootMode
	counter   float64

	bullets []*Bullet
}

func (p *Player) physics(dt float64) {
	p.rigidBody.physics(dt)
}

func (p *Player) upadte(dt float64, ovnis []*Ovni) []*Ovni {
	if top > 0 {
		p.rigidBody.velocity.Y += moveSpeed
	}

	if left > 0 {
		p.rigidBody.velocity.X -= moveSpeed
	}

	if right > 0 {
		p.rigidBody.velocity.X += moveSpeed
	}

	if bottom > 0 {
		p.rigidBody.velocity.Y -= moveSpeed
	}

	if enter > 0 {
		p.rigidBody.velocity = pixel.ZV
	}

	if space > 0 {
		if p.mode == shootLaser {
			x := (p.rigidBody.body.Min.X + p.rigidBody.body.Max.X) / 2
			rect := pixel.R(x, p.rigidBody.body.Min.Y, x+20, height)
			for _, o := range ovnis {
				if o.rigidBody.hit(rect) {
					o.loseLife(10)
				}
			}
		} else if p.mode == shootBullets && space%5 == 0 {
			p.bullets = append(p.bullets, &Bullet{
				rigidBody: NewRigidBodyBySize(p.rigidBody.body.Center().X,
					p.rigidBody.body.Center().Y,
					10, 10,
					pixel.V(math.Cos(p.counter*10), math.Sin(p.counter*10)).Scaled(100).Add(p.rigidBody.velocity),
				),
			})
		}
	}

	if tab == 1 {
		// works because only two modes
		p.mode = 1 - p.mode
		p.counter = 0
	}

	p.counter += dt

	p.physics(dt)

	var bullets []*Bullet
	for _, b := range p.bullets {
		r := b.rigidBody
		r.physics(dt)

		if r.body.Min.X < width/2 &&
			r.body.Min.Y < height &&
			r.body.Max.X > 0 &&
			r.body.Max.Y > 0 {

			hitted := false
			for _, o := range ovnis {
				if o.rigidBody.hit(b.rigidBody.body) {
					o.loseLife(5)
					hitted = true
				}
			}
			if !hitted {
				bullets = append(bullets, b)
			}
		}
	}

	var os []*Ovni
	for _, o := range ovnis {
		if o.isAlive() {
			os = append(os, o)
		}
	}

	ovnis = os

	p.bullets = bullets

	// TODO: bounce effect
	if p.rigidBody.body.Min.X < 0 {
		p.rigidBody.body = p.rigidBody.body.Moved(pixel.V(-p.rigidBody.body.Min.X, 0))
		p.rigidBody.velocity.X = 0
	}

	half := width / 2
	if p.rigidBody.body.Max.X > half {
		p.rigidBody.body = p.rigidBody.body.Moved(pixel.V(half-p.rigidBody.body.Max.X, 0))
		p.rigidBody.velocity.X = 0
	}

	if p.rigidBody.body.Max.Y < 0 {
		p.rigidBody.body = p.rigidBody.body.Moved(pixel.V(-p.rigidBody.body.Max.Y, 0))
		p.rigidBody.velocity.Y = 0
	}

	if p.rigidBody.body.Min.Y > height {
		p.rigidBody.body = p.rigidBody.body.Moved(pixel.V(height-p.rigidBody.body.Min.Y, 0))
		p.rigidBody.velocity.Y = 0
	}

	return ovnis
}

func (p *Player) draw(imag *imdraw.IMDraw) {
	imag.Color = colornames.Red
	p.rigidBody.draw(imag)

	imag.Color = colornames.Chartreuse
	if space > 0 && p.mode == shootLaser {
		x := (p.rigidBody.body.Min.X + p.rigidBody.body.Max.X) / 2
		imag.Push(pixel.V(x, height))
		imag.Push(pixel.V(x+20, p.rigidBody.body.Min.Y))
		imag.Rectangle(0)
	}

	for _, b := range p.bullets {
		b.rigidBody.draw(imag)
	}
}
