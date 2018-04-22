package main

import (
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

const (
	pauseMode     = 0
	asteroidField = 1
)

var (
	genMode    = asteroidField
	genCounter = 0
)

type kindOfOvni int

type Ovni struct {
	rigidBody *RigidBody
	life      int
}

func (o *Ovni) draw(imag *imdraw.IMDraw) {
	imag.Color = colornames.Brown
	o.rigidBody.draw(imag)
}

func (o *Ovni) loseLife(life int) {
	o.life -= life
	if life < 0 {
		o.life = 0
	}
}

func (o *Ovni) isAlive() bool {
	return o.life > 0
}

func updateUniverse(dt float64, ovnis []*Ovni) []*Ovni {
	var os []*Ovni
	for _, o := range ovnis {
		r := o.rigidBody
		r.physics(dt)
		if r.body.Max.Y > 0 {
			os = append(os, o)
		}
	}

	ovnis = os

	if genCounter%500 == 0 {
		genMode = (genMode + 1) % 3
		genCounter = 0
	}

	genCounter++

	switch genMode {
	case asteroidField:
		ovnis = genAsteroids(ovnis)
	}

	return ovnis
}

func genAsteroids(ovnis []*Ovni) []*Ovni {
	if genCounter%50 == 0 {
		alpha := 1 + rand.Float64()
		s := 20 * alpha
		ovnis = append(ovnis, &Ovni{
			rigidBody: NewRigidBodyBySize(
				rand.Float64()*width/2,
				height+50,
				s,
				s,
				pixel.V(0, -100*alpha),
			),
			life: int(2 * alpha),
		})
	}
	return ovnis
}
