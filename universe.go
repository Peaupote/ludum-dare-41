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
	piratesAttack = 2

	asteroid = 0
	pirates  = 1
)

var (
	genMode    = piratesAttack
	genCounter = 0
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

	if genCounter%1000 == 0 {
		genMode = (genMode + 1) % 3
		genCounter = 0
	}

	genCounter++

	switch genMode {
	case asteroidField:
		ovnis = genAsteroids(ovnis)
	case piratesAttack:
		ovnis = genPirates(ovnis)
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
			kind: asteroid,
		})
	}
	return ovnis
}

func genPirates(ovnis []*Ovni) []*Ovni {
	return ovnis
}
