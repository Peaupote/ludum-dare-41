package main

import (
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

const (
	asteroid = 0
)

type kindOfOvni int

type Ovni struct {
	rigidBody *RigidBody
	life      int
	kind      kindOfOvni
}

func (o *Ovni) draw(imag *imdraw.IMDraw) {
	imag.Color = colornames.Brown
	o.rigidBody.draw(imag)
}

func (o *Ovni) loseLife(life int) {
	life = o.life - life
	if life < 0 {
		o.life = 0
	} else {
		o.life = life
	}
}

func (o *Ovni) isAlive() bool {
	return o.life > 0
}

func updateUniverse(dt float64, ovnis []*Ovni) []*Ovni {
	for i, o := range ovnis {
		r := o.rigidBody
		r.physics(dt)
		if r.body.Max.Y < 0 {
			ovnis[i] = ovnis[len(ovnis)-1]
			ovnis = ovnis[:len(ovnis)-1]
		}
	}

	// TODO: less "random" random generation
	if rand.Intn(100) == 1 {
		ovnis = append(ovnis, &Ovni{
			rigidBody: NewRigidBodyBySize(
				rand.Float64()*(width/2),
				height+10,
				50.0, 50.0,
				pixel.V(0, -100),
			),
			life: rand.Intn(10) + 5,
			kind: asteroid,
		})
	}

	return ovnis
}
